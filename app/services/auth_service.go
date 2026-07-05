package services

import (
	"context"
	"errors"
	"strings"

	helper "app/helpers"
	model "app/models"
	repository "app/repositories"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (model.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}

	created, err := s.userRepo.Create(ctx, model.User{
		Name:         strings.TrimSpace(name),
		Email:        strings.TrimSpace(strings.ToLower(email)),
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		if isDuplicateEntryError(err) {
			return model.User{}, ErrEmailAlreadyExists
		}
		return model.User{}, err
	}

	return created, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (model.User, string, error) {
	user, err := s.userRepo.GetByEmail(ctx, strings.TrimSpace(strings.ToLower(email)))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return model.User{}, "", ErrInvalidCredentials
		}
		return model.User{}, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return model.User{}, "", ErrInvalidCredentials
	}

	token, err := helper.GenerateToken(user.ID, user.Email)
	if err != nil {
		return model.User{}, "", err
	}

	return user, token, nil
}

func isDuplicateEntryError(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "duplicate entry")
}
