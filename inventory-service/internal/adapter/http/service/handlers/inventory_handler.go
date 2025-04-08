package handlers

import (
	"inventory-service/internal/adapter/http/service/handlers/dto"
	"inventory-service/pkg/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Inventory struct {
	invUseCase InventoryUsecase
}

func NewInventory(invUseCase InventoryUsecase) *Inventory {
	return &Inventory{invUseCase: invUseCase}
}

func (h *Inventory) Create(ctx *gin.Context) {
	inventory, err := dto.ToInventoryCreateRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	v := validator.New()
	if dto.ValidateInventory(v, inventory); !v.Valid() {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": v.Errors})
		return
	}

	inventoryId, err := h.invUseCase.CreateItem(ctx.Request.Context(), inventory)
	if err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	inventory.ID = inventoryId
	ctx.JSON(http.StatusCreated, gin.H{"inventory": dto.ToInventoryCreateResponse(inventory)})

}

func (h *Inventory) Get(ctx *gin.Context) {

}

func (h *Inventory) GetListInventory(ctx *gin.Context) {

}

func (h *Inventory) Update(ctx *gin.Context) {

}

func (h *Inventory) Delete(ctx *gin.Context) {

}
