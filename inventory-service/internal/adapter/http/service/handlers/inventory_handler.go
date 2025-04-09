package handlers

import (
	"github.com/gin-gonic/gin"
)

type Inventory struct {
	invUseCase InventoryUsecase
}

func NewInventory(invUseCase InventoryUsecase) *Inventory {
	return &Inventory{invUseCase: invUseCase}
}

func (h *Inventory) Create(ctx *gin.Context) {

}

func (h *Inventory) Get(ctx *gin.Context) {

}

func (h *Inventory) GetListInventory(ctx *gin.Context) {

}

func (h *Inventory) Update(ctx *gin.Context) {

}

func (h *Inventory) Delete(ctx *gin.Context) {

}
