package invdto

import "order-service/internal/models"

// Inventory represents the inventory item structure
type Inventory struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Available   int64  `json:"available"`
	CreatedAt   string `json:"created_at"`
	Version     int32  `json:"version"`
}

// InventoryResponse represents the expected API response structure
type InventoryResponse struct {
	Inventory Inventory `json:"inventory"`
}

func ToInventoryModel(resp InventoryResponse) models.Inventory {
	return models.Inventory{
		ID:          resp.Inventory.ID,
		Name:        resp.Inventory.Name,
		Description: resp.Inventory.Description,
		Price:       resp.Inventory.Price,
		Available:   resp.Inventory.Available,
		CreatedAt:   resp.Inventory.CreatedAt,
		Version:     resp.Inventory.Version,
	}
}
