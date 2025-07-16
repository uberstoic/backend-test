package service

import (
	"context"

	"github.com/Serhio/backend-vk-test/internal/models"
	"github.com/Serhio/backend-vk-test/internal/storage/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User interface {
	Register(ctx context.Context, login, password string) (*models.User, error)
	Login(ctx context.Context, login, password string) (string, error)
}

type Ad interface {
	CreateAd(ctx context.Context, ad *models.Ad) (*models.Ad, error)
	GetAds(ctx context.Context, page, limit int, sortBy, sortDir string, minPrice, maxPrice float64, currentUserID int64) ([]models.Ad, error)
}

type Auth interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
	GenerateToken(user *models.User) (string, error)
	ParseToken(tokenString string) (int64, error)
}

// key struct of project
type Service struct {
	User User
	Ad   Ad
	Auth Auth
}

// constructor
func NewService(pool *pgxpool.Pool, jwtSecret string) *Service {
	userStorage := postgres.NewUserStorage(pool)
	adStorage := postgres.NewAdStorage(pool)

	authService := NewAuthService(jwtSecret)

	return &Service{
		User: NewUserService(userStorage, authService),
		Ad:   NewAdService(adStorage),
		Auth: authService,
	}
}
