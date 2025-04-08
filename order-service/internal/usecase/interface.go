package usecase

import (
	"context"
	"order-service/internal/models"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order models.Order) (int64, error)
	GetByFilter(ctx context.Context, filter models.OrderFilter) (models.Order, error)
	GetListOrderByFilter(ctx context.Context, filter models.OrderFilter) ([]models.Order, error)
	UpdateOrder(ctx context.Context, update models.OrderUpdateData) error

}
