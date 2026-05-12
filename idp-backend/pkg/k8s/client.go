package k8s

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewK8sClient(kubeconfigPath string, hostOverride string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Fatalf("Failed to build kubeconfig: %v", err)
	}

	// Override host if provided (useful for Docker-to-Host communication)
	if hostOverride != "" {
		config.Host = hostOverride
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create kubernetes client: %v", err)
	}

	log.Println("Kubernetes client initialized successfully")
	return clientset
}
