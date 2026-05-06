package service

import (
	"context"
	"fmt"
	"log"

	"idp-backend/internal/domain"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/ptr"
	"strings"
	"time"
)

type K8sService struct {
	client *kubernetes.Clientset
}

func NewK8sService(client *kubernetes.Clientset) *K8sService {
	return &K8sService{client: client}
}

func (s *K8sService) DeployApp(ctx context.Context, project *domain.Project, imageName string, port int) (string, error) {
	namespace := "idp-apps"
	projectName := project.Name
	finalURL := fmt.Sprintf("http://%s", project.Domain)

	// 1. Create Namespace if not exists
	_, _ = s.client.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: namespace},
	}, metav1.CreateOptions{})

	// 1.5 Delete existing resources if they exist (for clean redeploy)
	_ = s.client.AppsV1().Deployments(namespace).Delete(ctx, projectName, metav1.DeleteOptions{})
	_ = s.client.CoreV1().Services(namespace).Delete(ctx, projectName, metav1.DeleteOptions{})
	_ = s.client.NetworkingV1().Ingresses(namespace).Delete(ctx, projectName, metav1.DeleteOptions{})

	// 2. Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: projectName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": projectName},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": projectName},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  projectName,
							Image: imageName,
							Ports: []corev1.ContainerPort{{ContainerPort: int32(port)}},
						},
					},
				},
			},
		},
	}

	_, err := s.client.AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create deployment: %v", err)
	}

	// 3. Service
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: projectName,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": projectName},
			Ports: []corev1.ServicePort{
				{
					Port:       int32(port),
					TargetPort: intstr.FromInt(port),
				},
			},
		},
	}
	_, err = s.client.CoreV1().Services(namespace).Create(ctx, svc, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create service: %v", err)
	}

	// 4. Ingress
	pathType := networkingv1.PathTypePrefix
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: projectName,
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: ptr.To("nginx"),
			Rules: []networkingv1.IngressRule{
				{
					Host: project.Domain,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: projectName,
											Port: networkingv1.ServiceBackendPort{
												Number: int32(port),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	_, err = s.client.NetworkingV1().Ingresses(namespace).Create(ctx, ingress, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create ingress: %v", err)
	}

	// 5. Wait for Pod to be Running
	log.Printf("[K8S] Waiting for pod %s to be ready...", projectName)
	for i := 0; i < 20; i++ { // Wait up to 100 seconds
		time.Sleep(5 * time.Second)
		pods, err := s.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
			LabelSelector: fmt.Sprintf("app=%s", projectName),
		})
		if err != nil {
			continue
		}

		if len(pods.Items) > 0 {
			pod := pods.Items[0]
			if pod.Status.Phase == corev1.PodRunning {
				log.Printf("[K8S] Pod %s is now Running", projectName)
				return finalURL, nil
			}
			if pod.Status.Phase == corev1.PodFailed {
				return "", fmt.Errorf("pod failed to start")
			}
		}
	}

	return "", fmt.Errorf("timeout waiting for pod to be ready")
}

func (s *K8sService) GetLogs(ctx context.Context, projectName string) (string, error) {
	namespace := "idp-apps"

	// Find pods for this project
	pods, err := s.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", projectName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to list pods: %v", err)
	}

	if len(pods.Items) == 0 {
		return "No pods found for this project", nil
	}

	// Get logs from the first pod
	podName := pods.Items[0].Name
	logOptions := &corev1.PodLogOptions{}
	req := s.client.CoreV1().Pods(namespace).GetLogs(podName, logOptions)
	logs, err := req.DoRaw(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %v", err)
	}

	return string(logs), nil
}

func (s *K8sService) CreateBuildJob(ctx context.Context, project *domain.Project, dockerfile string, imageName string) error {
	namespace := "idp-apps"
	jobName := fmt.Sprintf("build-%s-%d", project.Name, time.Now().Unix())

	// 1. Create ConfigMap for Dockerfile
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
		},
		Data: map[string]string{
			"Dockerfile": dockerfile,
		},
	}
	_, err := s.client.CoreV1().ConfigMaps(namespace).Create(ctx, cm, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create build configmap: %v", err)
	}
	defer s.client.CoreV1().ConfigMaps(namespace).Delete(ctx, jobName, metav1.DeleteOptions{})

	// 2. Create Job
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "kaniko",
							Image: "gcr.io/kaniko-project/executor:latest",
							Args: []string{
								"--dockerfile=/workspace/Dockerfile",
								"--context=" + s.formatGitContext(project.GitURL),
								"--destination=" + imageName,
								"--insecure", // Assuming local registry without TLS
								"--skip-tls-verify",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "dockerfile",
									MountPath: "/workspace",
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
					Volumes: []corev1.Volume{
						{
							Name: "dockerfile",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: jobName,
									},
								},
							},
						},
					},
				},
			},
			BackoffLimit: ptr.To(int32(0)),
		},
	}

	_, err = s.client.BatchV1().Jobs(namespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create build job: %v", err)
	}

	// 3. Wait for completion
	for {
		time.Sleep(5 * time.Second)
		j, err := s.client.BatchV1().Jobs(namespace).Get(ctx, jobName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if j.Status.Succeeded > 0 {
			break
		}
		if j.Status.Failed > 0 {
			return fmt.Errorf("build job failed")
		}
	}

	return nil
}

func (s *K8sService) formatGitContext(gitURL string) string {
	// Kaniko expects git://github.com/user/repo.git format
	// If it starts with https:// or http://, strip it and replace with git://
	url := gitURL
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	
	if !strings.HasPrefix(url, "git://") {
		url = "git://" + url
	}
	return url
}

func int32Ptr(i int32) *int32 { return &i }
