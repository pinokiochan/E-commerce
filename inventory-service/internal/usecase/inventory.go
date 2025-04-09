package usecase

import (
	"context"
	"inventory-service/internal/adapter/http/service/handlers/dto"
	"inventory-service/internal/models"
	"inventory-service/pkg/validator"
)

type Inventory struct {
	invRepo InventoryRepository
}

func NewInventory(invRepo InventoryRepository) *Inventory {
	return &Inventory{invRepo: invRepo}
}

func (c *Inventory) CreateItem(ctx context.Context, request models.Inventory) (models.Inventory, error) {
	id, err := c.invRepo.CreateItem(ctx, request)
	if err != nil {
		return models.Inventory{}, err
	}
	request.ID = id
	return request, nil

}

func (c *Inventory) Get(ctx context.Context, id int64) (models.Inventory, error) {
	inv, err := c.invRepo.Get(ctx, id)

	if err != nil {
		return models.Inventory{}, err
	}

	return inv, nil
}

func (c *Inventory) GetListInventory(ctx context.Context, filters models.Filters) ([]models.Inventory, dto.Metadata, error) {
	v := validator.New()
	models.ValidateFilters(v, filters)

	if !v.Valid() {
		return nil, dto.Metadata{}, dto.ErrInvalidFilters
	}

	items, totalRecords, err := c.invRepo.GetListInventory(ctx, filters)

	if err != nil {
		return nil, dto.Metadata{}, err
	}

	metadata := dto.CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return items, metadata, nil
}

func (c *Inventory) Update(ctx context.Context, request models.UpdateInventoryData) (models.Inventory, error) {
	item, err := c.invRepo.Get(ctx, *request.ID)
	if err != nil {
		return models.Inventory{}, err
	}

	if request.Version != nil && item.Version != *request.Version {
		return models.Inventory{}, dto.ErrEditConflict
	}

	// If the input.Name value is nil then we know that no corresponding "name" key/
	// value pair was provided in the request body. So we move on and leave the
	// movie record unchanged. Otherwise, we update the movie record with the new name
	// value. Importantly, because input.Name is a now a pointer to a string, we need
	// to dereference the pointer using the * operator to get the underlying value
	// before assigning it to our movie record.
	if request.Name != nil {
		item.Name = *request.Name
	}
	if request.Description != nil {
		item.Description = *request.Description
	}
	if request.Price != nil {
		item.Price = *request.Price
	}
	if request.Available != nil {
		item.Available = *request.Available
	}

	v := validator.New()
	if dto.ValidateInventory(v, item); !v.Valid() {
		return item, dto.ErrUnprocessableEntity
	}

	err = c.invRepo.Update(ctx, &item)
	if err != nil {
		return models.Inventory{}, err
	}

	return item, nil
}

func (c *Inventory) Delete(ctx context.Context, id int64) error {
	err := c.invRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
