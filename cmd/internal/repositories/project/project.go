package project

import "prodyo-backend/cmd/internal/models"

type Project struct {
	projects []models.Project
}

func New() *Project {
	return &Project{}
}

func (p Project) GetAll() ([]models.Project, error) {
	return p.projects, error(nil)
}

func (p *Project) Add(newProject models.Project) error {
	p.projects = append(p.projects, newProject)
	return nil
}
