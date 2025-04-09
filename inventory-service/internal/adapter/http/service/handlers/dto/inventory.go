package dto

import (
	"inventory-service/internal/adapter/postgres/dao"
	"inventory-service/internal/models"
	"inventory-service/pkg/validator"
	"strconv"
	"time"

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

type InventoryUpdateRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Available   *int64   `json:"available"`
}

type InventoryResponse struct {
	ID          int64     `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price,omitempty"`
	Available   int64     `json:"available,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitzero"`
	Version     int32     `json:"version,omitempty"`
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

func ToInventoryUpdateRequest(ctx *gin.Context) (models.UpdateInventoryData, error) {
	id, err := ReadParamID(ctx)
	if err != nil {
		return models.UpdateInventoryData{}, err
	}

	var req InventoryUpdateRequest
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return models.UpdateInventoryData{}, err
	}

	var inventory models.UpdateInventoryData

	expectedVersion := ctx.GetHeader("X-Expected-Version")

	if expectedVersion != "" {
		expectedVersionInt, err := strconv.ParseInt(expectedVersion, 10, 32)
		if err != nil {
			return models.UpdateInventoryData{}, err
		}
		result := int32(expectedVersionInt)
		inventory.Version = &result
	}

	inventory.ID = &id
	inventory.Name = req.Name
	inventory.Description = req.Description
	inventory.Price = req.Price
	inventory.Available = req.Available

	return inventory, nil
}

func ToInventoryCreateResponse(inv models.Inventory) InventoryCreateResponse {
	return InventoryCreateResponse{
		ID:   inv.ID,
		Name: inv.Name,
	}
}

func ToInventoryListResponse(invs []models.Inventory) []InventoryResponse {
	var responce []InventoryResponse

	for _, v := range invs {
		responce = append(responce, ToInventoryResponse(v))
	}

	return responce
}

func ToInventoryResponse(inv models.Inventory) InventoryResponse {
	return InventoryResponse{
		ID:          inv.ID,
		Name:        inv.Name,
		Description: inv.Description,
		Price:       inv.Price,
		Available:   inv.Available,
		CreatedAt:   inv.CreatedAt,
		Version:     inv.Version,
	}
}

func ParseListRequest(ctx *gin.Context, v *validator.Validator) models.Filters {
	filter := models.Filters{
		Page:         1,                // Default current page
		PageSize:     8,                // Default page size
		Sort:         "id",             // Default sort value
		SortSafelist: dao.SafeSortList, // Available sort options
	}

	// Parse page parameter
	filter.Page = ReadInt(ctx, "page", filter.Page, v)

	// Parse page_size parameter
	filter.PageSize = ReadInt(ctx, "page_size", filter.PageSize, v)

	// Parse sort parameter
	filter.Sort = ReadString(ctx, "sort", filter.Sort)

	if sortParam := ctx.Query("sort"); sortParam != "" {
		filter.Sort = sortParam
	}

	models.ValidateFilters(v, filter)

	return filter
}
