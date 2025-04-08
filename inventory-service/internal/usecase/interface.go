package usecase

import (
	"api-gateway/internal/models"
	"context"
)

type InventoryRepository interface {
	CreateItem(ctx context.Context, item models.Inventory) (int64, error)
	Get(ctx context.Context, id int64) (models.Inventory, error)
	GetListInventory(ctx context.Context, filters models.Filters) ([]models.Inventory, error)
	Update(ctx context.Context, item models.Inventory) error
	Delete(ctx context.Context, id int64) error
}
