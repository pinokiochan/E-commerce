package handlers

import (
	"errors"
	"inventory-service/internal/adapter/http/service/handlers/dto"
	"inventory-service/pkg/validator"
	"log"
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

	inventoryNew, err := h.invUseCase.CreateItem(ctx.Request.Context(), inventory)
	if err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"inventory": dto.ToInventoryCreateResponse(inventoryNew)})
}

func (h *Inventory) GetList(ctx *gin.Context) {
	v := validator.New()

	filters := dto.ParseListRequest(ctx, v)
	if !v.Valid() {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": v.Errors})
		return
	}

	items, metadata, err := h.invUseCase.GetListInventory(ctx.Request.Context(), filters)
	if err != nil {
		errCtx := dto.FromError(err)
		log.Println(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"inventory": dto.ToInventoryListResponse(items),
		"metadata":  metadata,
	})
}

func (h *Inventory) GetByID(ctx *gin.Context) {
	id, err := dto.ReadParamID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	inventory, err := h.invUseCase.Get(ctx.Request.Context(), id)
	if err != nil {
		log.Println(err)
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"inventory": dto.ToInventoryResponse(inventory)})
}

func (h *Inventory) Update(ctx *gin.Context) {
	item, err := dto.ToInventoryUpdateRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemUpdated, err := h.invUseCase.Update(ctx.Request.Context(), item)
	if err != nil {
		if errors.Is(err, dto.ErrUnprocessableEntity) {
			v := validator.New()
			dto.ValidateInventory(v, itemUpdated)
			if !v.Valid() {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": v.Errors})
				return
			}
		}
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"inventory": dto.ToInventoryResponse(itemUpdated)})
}

func (h *Inventory) Delete(ctx *gin.Context) {
	id, err := dto.ReadParamID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.invUseCase.Delete(ctx.Request.Context(), id)
	if err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	ctx.Status(http.StatusNoContent)
}
