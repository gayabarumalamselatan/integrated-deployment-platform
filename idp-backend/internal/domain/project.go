package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProjectStatus string

const (
	StatusPending  ProjectStatus = "PENDING"
	StatusBuilding ProjectStatus = "BUILDING"
	StatusDeploying ProjectStatus = "DEPLOYING"
	StatusReady    ProjectStatus = "READY"
	StatusFailed   ProjectStatus = "FAILED"
)

type Project struct {
	ID        uuid.UUID     `db:"id" json:"id"`
	Name      string        `db:"name" json:"name"`
	GitURL    string        `db:"git_url" json:"git_url"`
	Domain    string        `db:"domain" json:"domain"`
	Status    ProjectStatus `db:"status" json:"status"`
	CreatedAt time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt time.Time     `db:"updated_at" json:"updated_at"`
}

type ProjectRepository interface {
	Create(project *Project) error
	GetAll() ([]Project, error)
	GetByID(id uuid.UUID) (*Project, error)
	UpdateStatus(id uuid.UUID, status ProjectStatus) error
}

type ProjectService interface {
	CreateProject(name, gitURL string) (*Project, error)
	GetProjects() ([]Project, error)
	GetProjectByID(id uuid.UUID) (*Project, error)
	GetProjectLogs(id uuid.UUID) (string, error)
	ProcessDeployment(projectID uuid.UUID) error
}
