package usecases

import (
	"context"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/repositories/user"

	"github.com/google/uuid"
)

type UserUseCase struct {
	repo *user.Repository
}

// Constructor
func NewUserUseCase(repo *user.Repository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (u *UserUseCase) GetAll(ctx context.Context, pagination models.PaginationRequest) ([]models.User, models.PaginationResponse, error) {
	return u.repo.GetAll(ctx, pagination)
}

func (u *UserUseCase) GetByID(ctx context.Context, id uuid.UUID) (models.User, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *UserUseCase) GetByEmail(ctx context.Context, email string) (models.User, error) {
	return u.repo.GetByEmail(ctx, email)
}

func (u *UserUseCase) Add(ctx context.Context, newUser models.User) (uuid.UUID, error) {
	if newUser.ID == uuid.Nil {
		newUser.ID = uuid.New()
	}

	if err := u.repo.Add(ctx, newUser); err != nil {
		return uuid.Nil, err
	}

	return newUser.ID, nil
}

func (u *UserUseCase) Update(ctx context.Context, user models.User) error {
	return u.repo.Update(ctx, user)
}

func (u *UserUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
