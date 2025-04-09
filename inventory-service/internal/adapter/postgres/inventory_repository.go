package postgres

import (
	"context"
	"fmt"
	"inventory-service/internal/adapter/postgres/dao"
	"inventory-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryRepository struct {
	db *pgxpool.Pool
}

func NewInventoryRepository(db *pgxpool.Pool) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (p *InventoryRepository) CreateItem(ctx context.Context, item models.Inventory) (int64, error) {
	query := `
		INSERT INTO inventory (name, description, price, available)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int64
	err := p.db.QueryRow(ctx, query,
		item.Name,
		item.Description,
		item.Price,
		item.Available,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil

}

func (p *InventoryRepository) Get(ctx context.Context, id int64) (models.Inventory, error) {
	query := `
		SELECT id, created_at, name, description, price, available, isdeleted, version
		from inventory
		WHERE id = $1 AND isdeleted = false
	`
	var item models.Inventory
	err := p.db.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.CreatedAt,
		&item.Name,
		&item.Description,
		&item.Price,
		&item.Available,
		&item.IsDeleted,
		&item.Version,
	)
	if err != nil {
		return models.Inventory{}, err
	}
	return item, nil
}

func (p *InventoryRepository) GetListInventory(ctx context.Context, filters models.Filters) ([]models.Inventory, int, error) {
	query := fmt.Sprintf(`
	SELECT count(*) OVER(), id, created_at, name, description, price, available, isdeleted, version
	FROM inventory
	WHERE isdeleted = false
	ORDER BY %s %s, id ASC
	LIMIT $1 OFFSET $2`, filters.SortColumn(), filters.SortDirection())

	args := []any{filters.Limit(), filters.Offset()}

	var total_records int
	var items []models.Inventory

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return []models.Inventory{}, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var item models.Inventory
		err := rows.Scan(
			&total_records,
			&item.ID,
			&item.CreatedAt,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.Available,
			&item.IsDeleted,
			&item.Version,
		)

		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total_records, nil

}

func (p *InventoryRepository) Update(ctx context.Context, item *models.Inventory) error {
	query := `
		UPDATE inventory
		SET name = $1, description = $2, price = $3, available = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version
	`
	args := []any{
		item.Name,
		item.Description,
		item.Price,
		item.Available,
		item.ID,
		item.Version,
	}

	err := p.db.QueryRow(ctx, query, args...).Scan(&item.Version)
	if err != nil {
		return err
	}

	return nil
}

func (p *InventoryRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE inventory 
		SET isdeleted = true,
			version= version + 1
		WHERE id = $1
	`

	result, err := p.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return dao.ErrRecordNotFound
	}

	return nil
}
