package project

import (
	"errors"
	"prodyo-backend/cmd/internal/models"

	"github.com/google/uuid"
)

type Project struct {
	projects []models.Project
}

func New() *Project {
	return &Project{}
}

func (p Project) GetAll() ([]models.Project, error) {
	return p.projects, nil
}

func (p Project) GetByID(id uuid.UUID) (models.Project, error) {
	for _, project := range p.projects {
		if project.ID == id {
			return project, nil
		}
	}
	return models.Project{}, errors.New("project not found")
}

func (p *Project) Add(newProject models.Project) error {
	p.projects = append(p.projects, newProject)
	return nil
}

func (p *Project) Update(project models.Project) error {
	for i, existingProject := range p.projects {
		if existingProject.ID == project.ID {
			p.projects[i] = project
			return nil
		}
	}
	return errors.New("project not found")
}

func (p *Project) Delete(id uuid.UUID) error {
	for i, project := range p.projects {
		if project.ID == id {
			p.projects = append(p.projects[:i], p.projects[i+1:]...)
			return nil
		}
	}
	return errors.New("project not found")
}
