package service

import (
	"context"
	"fmt"
	"log"

	"idp-backend/internal/domain"

	"github.com/google/uuid"
)

type projectService struct {
	repo         domain.ProjectRepository
	kafkaService *KafkaService
	k8sService   *K8sService
	buildService *BuildService
	baseDomain   string
}

func NewProjectService(
	repo domain.ProjectRepository,
	kafkaService *KafkaService,
	k8sService *K8sService,
	buildService *BuildService,
	baseDomain string,
) domain.ProjectService {
	return &projectService{
		repo:         repo,
		kafkaService: kafkaService,
		k8sService:   k8sService,
		buildService: buildService,
		baseDomain:   baseDomain,
	}
}

func (s *projectService) CreateProject(name, gitURL string) (*domain.Project, error) {
	project := &domain.Project{
		ID:     uuid.New(),
		Name:   name,
		GitURL: gitURL,
		Domain: fmt.Sprintf("%s.%s", name, s.baseDomain),
		Status: domain.StatusPending,
	}

	if err := s.repo.Create(project); err != nil {
		return nil, fmt.Errorf("failed to save project: %v", err)
	}

	// Push to Kafka
	if err := s.kafkaService.ProduceDeployRequest(project.ID); err != nil {
		log.Printf("Warning: Failed to produce kafka message for %s: %v", project.ID, err)
		// We still return success since it's saved in DB, but background worker might miss it
	}

	return project, nil
}

func (s *projectService) GetProjects() ([]domain.Project, error) {
	return s.repo.GetAll()
}

func (s *projectService) GetProjectByID(id uuid.UUID) (*domain.Project, error) {
	return s.repo.GetByID(id)
}

func (s *projectService) GetProjectLogs(id uuid.UUID) (string, error) {
	project, err := s.repo.GetByID(id)
	if err != nil {
		return "", err
	}

	return s.k8sService.GetLogs(context.Background(), project.Name)
}

func (s *projectService) ProcessDeployment(projectID uuid.UUID) error {
	project, err := s.repo.GetByID(projectID)
	if err != nil {
		return err
	}

	// Update status to BUILDING
	s.repo.UpdateStatus(projectID, domain.StatusBuilding)

	// 1. Build Image
	buildResult, err := s.buildService.Build(project)
	if err != nil {
		log.Printf("Build failed for %s: %v", project.Name, err)
		s.repo.UpdateStatus(projectID, domain.StatusFailed)
		return err
	}

	// 2. Trigger K8s deployment
	s.repo.UpdateStatus(projectID, domain.StatusDeploying)
	finalURL, err := s.k8sService.DeployApp(context.Background(), project, buildResult.ImageName, buildResult.Port)
	if err != nil {
		log.Printf("Deployment failed for %s: %v", project.Name, err)
		s.repo.UpdateStatus(projectID, domain.StatusFailed)
		return err
	}

	// Update domain/URL and status to READY
	s.repo.UpdateDomain(projectID, finalURL)
	s.repo.UpdateStatus(projectID, domain.StatusReady)
	return nil
}
