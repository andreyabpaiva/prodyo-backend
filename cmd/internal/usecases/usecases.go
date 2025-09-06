package usecases

import "prodyo-backend/cmd/internal/repositories"

type ProjectUseCase struct {
	repos repositories.Repositories
}

func New(repos *repositories.Repositories) *ProjectUseCase {
	return &ProjectUseCase{
		repos: *repos,
	}
}
