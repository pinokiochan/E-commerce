package handlers

import (
	"context"
	"inventory-service/internal/adapter/http/service/handlers/dto"
	"inventory-service/internal/models"
)

type InventoryUsecase interface {
	CreateItem(ctx context.Context, request models.Inventory) (models.Inventory, error)
	Get(ctx context.Context, id int64) (models.Inventory, error)
	GetListInventory(ctx context.Context, filters models.Filters) ([]models.Inventory, dto.Metadata, error)
	Update(ctx context.Context, request models.UpdateInventoryData) (models.Inventory, error)
	Delete(ctx context.Context, id int64) error
}
