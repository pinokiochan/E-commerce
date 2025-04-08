package dto

import (
	"inventory-service/internal/models"

	"github.com/gin-gonic/gin"
)

type InventoryCreateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Available   int64   `json:"available"`
}

type InventoryCreateResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func ToInventoryCreateRequest(ctx *gin.Context) (models.Inventory, error) {
	var req InventoryCreateRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		return models.Inventory{}, err
	}

	inventory := models.Inventory{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Available:   req.Available,
	}

	return inventory, nil
}

func ToInventoryCreateResponse(inv models.Inventory) InventoryCreateResponse {
	return InventoryCreateResponse{
		ID:   inv.ID,
		Name: inv.Name,
	}
}