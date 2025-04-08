package usecase

import (
	"api-gateway/internal/models"
	"context"
)

type Inventory struct {
	invRepo InventoryRepository
}

func NewInventory(invRepo InventoryRepository) *Inventory {
	return &Inventory{invRepo: invRepo}
}

func (c *Inventory) CreateItem(ctx context.Context, request models.Inventory) (models.Inventory, error) {
	panic("implement me")
}

func (c *Inventory) Get(ctx context.Context, id int64) (models.Inventory, error) {
	panic("implement me")
}

func (c *Inventory) GetListInventory(ctx context.Context, filters models.Filters) ([]models.Inventory, error) {
	panic("implement me")
}

func (c *Inventory) Update(ctx context.Context, item models.Inventory) error {
	panic("implement me")
}

func (c *Inventory) Delete(ctx context.Context, id int64) error {
	panic("implement me")
}
