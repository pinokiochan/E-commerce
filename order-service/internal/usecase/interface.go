package usecase

import (
	"context"
	"order-service/internal/models"
)

type OrderRepository interface {
	Create(ctx context.Context, order models.Order) (int64, error)
	GetWithFilter(ctx context.Context, filter models.OrderFilter) (models.Order, error)
	GetListWithFilter(ctx context.Context, filter models.OrderFilter) ([]models.Order, error)
	Update(ctx context.Context, update models.OrderUpdateData) error
}

type InventoryService interface {
	GetById(id int64) (models.Inventory, error)
	Substruct(id, newAvailability int64, version int32) error
}
