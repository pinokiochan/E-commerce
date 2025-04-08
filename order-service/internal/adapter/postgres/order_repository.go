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
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	query := `
			INSERT INTO orders (customername, status)
			VALUES ($1, $2)
			RETURNING ID;
	`
	var orderID int64
	err = tx.QueryRow(ctx, query, order.CustomerName, order.Status).Scan(&orderID)
	if err != nil {
		return 0, err
	}
	queryOrderItems := `
				INSERT INTO order_items (OrderID, ProductID, Quantity) VALUES
				($1, $2, $3)
				ON CONFLICT (OrderID, ProductID)
				DO UPDATE SET Quantity = order_items.Quantity + EXCLUDED.Quantity
	`
	for _, v := range order.OrderItems {
		_, err = tx.Exec(ctx, queryOrderItems, orderID, v.ProductID, v.Quantity)
		if err != nil {
			return 0, err
		}

	}
	return orderID, tx.Commit(ctx)

}
