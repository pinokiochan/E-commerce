package postgres

import (
	"inventory-service/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryRepository struct {
	db *pgxpool.Pool
}

func NewInventoryRepository(db *pgxpool.Pool) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (p *InventoryRepository) Create(ctx context.Context, item models.Inventory) (int64, error) {
	panic("Implement me")
}

func (p *InventoryRepository) Get(ctx context.Context, id int64) (models.Inventory, error) {
	panic("implement me")
}

func (p *InventoryRepository) GetListInventory(ctx context.Context, filters models.Filters)([]models.Inventory, error) {
	panic("implement me")
}

func (p *InventoryRepository) Update(ctx context.Context, item models.Inventory) error {
	panic("implement me")
}

func(p *InventoryRepository) Delete(ctx context.Context, id int64) error {
	panic("implement me")
}