package usecases

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/repositories/session"
	"prodyo-backend/cmd/internal/repositories/user"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrSessionExpired     = errors.New("session expired")
)

type AuthUseCase struct {
	userRepo    *user.Repository
	sessionRepo *session.Repository
}

func NewAuthUseCase(userRepo *user.Repository, sessionRepo *session.Repository) *AuthUseCase {
	return &AuthUseCase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (a *AuthUseCase) Register(ctx context.Context, email, password, name string) (uuid.UUID, error) {
	// Check if user already exists
	_, err := a.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return uuid.Nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}

	newUser := models.User{
		ID:           uuid.New(),
		Email:        email,
		Name:         name,
		PasswordHash: string(hashedPassword),
	}

	if err := a.userRepo.Add(ctx, newUser); err != nil {
		return uuid.Nil, err
	}

	return newUser.ID, nil
}

func (a *AuthUseCase) Login(ctx context.Context, email, password string) (models.Session, error) {
	user, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		log.Printf("Login: User not found for email %s: %v", email, err)
		return models.Session{}, ErrInvalidCredentials
	}

	if user.PasswordHash == "" {
		log.Printf("Login: User %s has no password_hash set", email)
		return models.Session{}, ErrInvalidCredentials
	}

	log.Printf("Login: Comparing password for user %s (hash length: %d)", email, len(user.PasswordHash))
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Printf("Login: Password comparison failed for user %s: %v", email, err)
		return models.Session{}, ErrInvalidCredentials
	}

	log.Printf("Login: Password verified successfully for user %s", email)

	token, err := generateToken()
	if err != nil {
		return models.Session{}, err
	}

	session := models.Session{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := a.sessionRepo.Create(ctx, session); err != nil {
		return models.Session{}, err
	}

	return session, nil
}

func (a *AuthUseCase) ValidateSession(ctx context.Context, token string) (models.User, error) {
	session, err := a.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		return models.User{}, ErrSessionExpired
	}

	user, err := a.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (a *AuthUseCase) Logout(ctx context.Context, token string) error {
	return a.sessionRepo.Delete(ctx, token)
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
