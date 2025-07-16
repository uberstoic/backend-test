package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Serhio/backend-vk-test/internal/models"
)

var (
	ErrTitleTooLong   = errors.New("title cannot be longer than 200 characters")
	ErrTextTooLong    = errors.New("text cannot be longer than 1000 characters")
	ErrPriceInvalid   = errors.New("price must be a non-negative number")
	ErrImageURLInvalid = errors.New("image URL is required")
)

type AdStorage interface {
	CreateAd(ctx context.Context, ad *models.Ad) (int64, error)
	GetAds(ctx context.Context, page, limit int, sortBy, sortDir string, minPrice, maxPrice float64, currentUserID int64) ([]models.Ad, error)
}

type AdService struct {
	storage AdStorage
}

func NewAdService(storage AdStorage) *AdService {
	return &AdService{storage: storage}
}

func (s *AdService) CreateAd(ctx context.Context, ad *models.Ad) (*models.Ad, error) {
	if err := validateAd(ad); err != nil {
		return nil, err
	}

	id, err := s.storage.CreateAd(ctx, ad)
	if err != nil {
		return nil, fmt.Errorf("failed to create ad: %w", err)
	}

	ad.ID = id
	return ad, nil
}

func (s *AdService) GetAds(ctx context.Context, page, limit int, sortBy, sortDir string, minPrice, maxPrice float64, currentUserID int64) ([]models.Ad, error) {
	return s.storage.GetAds(ctx, page, limit, sortBy, sortDir, minPrice, maxPrice, currentUserID)
}

func validateAd(ad *models.Ad) error {
	if len(ad.Title) > 200 {
		return ErrTitleTooLong
	}
	if len(ad.Text) > 1000 {
		return ErrTextTooLong
	}
	if ad.Price < 0 {
		return ErrPriceInvalid
	}
	if ad.ImageURL == "" {
		return ErrImageURLInvalid
	}
	return nil
}
