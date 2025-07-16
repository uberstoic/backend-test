package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/Serhio/backend-vk-test/internal/models"
	"github.com/Serhio/backend-vk-test/internal/storage/postgres"
)

var (
	ErrLoginFormat      = errors.New("login must be 3-30 characters long and contain only letters, numbers, and underscores")
	ErrPasswordFormat   = errors.New("password must be at least 6 characters long")
	ErrUserAlreadyExists = errors.New("user with this login already exists")
)

type UserStorage interface {
	CreateUser(ctx context.Context, user *models.User) (int64, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
}

type UserService struct {
	storage     UserStorage
	authService *AuthService
}

func NewUserService(storage UserStorage, authService *AuthService) *UserService {
	return &UserService{storage: storage, authService: authService}
}

func (s *UserService) Register(ctx context.Context, login, password string) (*models.User, error) {
	if err := validateCredentials(login, password); err != nil {
		return nil, err
	}

	_, err := s.storage.GetUserByLogin(ctx, login)
	if !errors.Is(err, postgres.ErrUserNotFound) {
		if err == nil {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	hash, err := s.authService.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Login:        login,
		PasswordHash: hash,
	}

	id, err := s.storage.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = id
	return user, nil
}

func (s *UserService) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.storage.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return "", postgres.ErrUserNotFound
		}
		return "", err
	}

	if !s.authService.CheckPasswordHash(password, user.PasswordHash) {
		return "", ErrInvalidPassword
	}

	return s.authService.GenerateToken(user)
}

func validateCredentials(login, password string) error {
	if len(password) < 6 {
		return ErrPasswordFormat
	}
	loginRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
	if !loginRegex.MatchString(login) {
		return ErrLoginFormat
	}
	return nil
}
