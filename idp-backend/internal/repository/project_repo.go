package repository

import (
	"idp-backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type projectRepo struct {
	db *sqlx.DB
}

func NewProjectRepository(db *sqlx.DB) domain.ProjectRepository {
	return &projectRepo{db: db}
}

func (r *projectRepo) Create(p *domain.Project) error {
	query := `INSERT INTO projects (id, name, git_url, domain, status, updated_at) 
              VALUES (:id, :name, :git_url, :domain, :status, NOW())`
	_, err := r.db.NamedExec(query, p)
	return err
}

func (r *projectRepo) GetAll() ([]domain.Project, error) {
	projects := []domain.Project{}
	err := r.db.Select(&projects, "SELECT * FROM projects ORDER BY created_at DESC")
	return projects, err
}

func (r *projectRepo) GetByID(id uuid.UUID) (*domain.Project, error) {
	var project domain.Project
	err := r.db.Get(&project, "SELECT * FROM projects WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepo) UpdateStatus(id uuid.UUID, status domain.ProjectStatus) error {
	_, err := r.db.Exec("UPDATE projects SET status = $1, updated_at = NOW() WHERE id = $2", status, id)
	return err
}

func (r *projectRepo) UpdateDomain(id uuid.UUID, domain string) error {
	_, err := r.db.Exec("UPDATE projects SET domain = $1, updated_at = NOW() WHERE id = $2", domain, id)
	return err
}
