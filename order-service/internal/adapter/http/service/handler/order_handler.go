package handler

import (
	"net/http"
	"order-service/internal/adapter/http/service/handler/dto"
	"order-service/internal/models"
	"order-service/pkg/validator"

	"github.com/gin-gonic/gin"
)

// OrderHandler
type Order struct {
	uc OrderUsecase
}

func NewOrder(uc OrderUsecase) *Order {
	return &Order{
		uc: uc,
	}
}

func (c *Order) Create(ctx *gin.Context) {
	order, err := dto.FromOrderCreateRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// fmt.Println(order.OrderItems[0].Quantity)
	v := validator.New()

	if dto.ValidateOrder(v, order); !v.Valid() {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": v.Errors})
		return
	}

	newOrder, err := c.uc.Create(ctx.Request.Context(), order)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"order": dto.ToOrderResponse(newOrder)})
}

func (c *Order) GetList(ctx *gin.Context) {
	orders, err := c.uc.GetList(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"orders": dto.ToOrderListResponse(orders)})
}

func (c *Order) GetByID(ctx *gin.Context) {
	id, err := dto.ReadParamID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	// Get order from service
	order, err := c.uc.Get(ctx.Request.Context(), id)
	if err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"order": dto.ToOrderResponse(order)})
}

func (c *Order) SetStatus(ctx *gin.Context) {
	id, err := dto.ReadParamID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	var request dto.OrderStatusRequest

	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}

	v := validator.New()
	if dto.ValidateSetOrderStatusRequest(v, request); !v.Valid() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": v.Errors})
		return
	}

	order, err := c.uc.SetStatus(ctx.Request.Context(), models.UpdateStatus{
		OrderID: id,
		Status:  request.Status,
	})
	if err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"order": dto.ToOrderResponse(order)})
}
