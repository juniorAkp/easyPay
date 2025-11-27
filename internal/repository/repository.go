package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/juniorAkp/easyPay/internal/model"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateUser(ctx context.Context, details *model.Details) error {
	_, err := r.pool.Exec(ctx,
		"INSERT INTO details(username, phone, updatedAt, createdAt) VALUES($1, $2, $3, $4)",
		details.Username, details.Phone, time.Now(), time.Now(),
	)
	return err
}

func (r *Repository) GetUserByPhone(ctx context.Context, phone string) (*model.Details, error) {
	var details model.Details

	err := r.pool.QueryRow(ctx, "SELECT id,username,phone FROM details WHERE phone = $1", phone).Scan(&details.ID, &details.Username, &details.Phone)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("no user found")
	} else if err != nil {
		return nil, err
	}

	return &details, nil
}
