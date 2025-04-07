package usecase

import (
	"context"
	"order-service/internal/models"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context, order models.Order) (int64, error)
	GetByFilter(ctx context.Context, filter models.OrderFilter) (models.Order, error)
	GetListOrderByFilter(ctx context.Context, filter models.OrderFilter) ([]models.Order, error)
	UpdateOrder(ctx context.Context, update models.OrderUpdateData) error
	DeleteOrder(ctx context.Context, id int64) error
}
