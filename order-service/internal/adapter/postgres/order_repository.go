package postgres

import (
	"context"
	"order-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Order struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *Order {
	return &Order{db: db}
}

func (r *Order) Create(ctx context.Context, order models.Order) (int64, error) {

}
