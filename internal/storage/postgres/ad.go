package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Serhio/backend-vk-test/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdStorage struct {
	pool *pgxpool.Pool
}

func NewAdStorage(pool *pgxpool.Pool) *AdStorage {
	return &AdStorage{pool: pool}
}

func (s *AdStorage) CreateAd(ctx context.Context, ad *models.Ad) (int64, error) {
	query := `INSERT INTO ads (title, text, image_url, price, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int64
	err := s.pool.QueryRow(ctx, query, ad.Title, ad.Text, ad.ImageURL, ad.Price, ad.UserID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *AdStorage) GetAds(ctx context.Context, page, limit int, sortBy, sortDir string, minPrice, maxPrice float64, currentUserID int64) ([]models.Ad, error) {
	query := `SELECT a.id, a.title, a.text, a.image_url, a.price, a.created_at, u.login, (a.user_id = $1) AS is_owner
	          FROM ads a JOIN users u ON a.user_id = u.id`

	var conditions []string
	var args []interface{}
	args = append(args, currentUserID)
	argID := 2 // start with $2 since $1 is currentUserID

	if minPrice > 0 {
		conditions = append(conditions, fmt.Sprintf("a.price >= $%d", argID))
		args = append(args, minPrice)
		argID++
	}
	if maxPrice > 0 {
		conditions = append(conditions, fmt.Sprintf("a.price <= $%d", argID))
		args = append(args, maxPrice)
		argID++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortDir)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, limit, (page-1)*limit)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ads []models.Ad
	for rows.Next() {
		var ad models.Ad
		if err := rows.Scan(&ad.ID, &ad.Title, &ad.Text, &ad.ImageURL, &ad.Price, &ad.CreatedAt, &ad.Author, &ad.IsOwner); err != nil {
			return nil, err
		}
		ads = append(ads, ad)
	}

	return ads, nil
}
