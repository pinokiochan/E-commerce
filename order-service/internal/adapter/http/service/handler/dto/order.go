package dto

import (
	"order-service/internal/adapter/postgres/dao"
	"order-service/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

type OrderCreateRequest struct {
	CustomerName string              `json:"customer_name"`
	OrderItems   []OrderItemsRequest `json:"items"`
}

type OrderItemsRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
}

type OrderCreateResponseRequest struct {
	OrderID      int    `json:"order_id"`
	CustomerName string `json:"customer_name"`
}
type OrderResponse struct {
	OrderID      int64               `json:"order_id"`
	CustomerName string              `json:"customer_name"`
	Items        []OrderItemsRequest `json:"items"`
	Status       string              `json:"status"`
	CreatedAt    time.Time           `json:"created_at"`
}

type OrderStatusRequest struct {
	Status string `json:"status"`
}

func FromOrderCreateRequest(ctx *gin.Context) (models.Order, error) {
	var req OrderCreateRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		return models.Order{}, err
	}
	var order models.Order
	order.CustomerName = req.CustomerName
	order.Status = dao.OrderStatusPending

	for _, v := range req.OrderItems {
		orderItems := models.OrderItem{
			ProductID: v.ProductID,
			Quantity:  v.Quantity,
		}
		order.OrderItems = append(order.OrderItems, orderItems)
	}
	return order, nil
}

func ToOrderListResponse(orders []models.Order) []OrderResponse {
	resp := []OrderResponse{}
	for _, order := range orders {
		orderResponse := ToOrderResponse(order)
		resp = append(resp, orderResponse)
	}
	return resp
}

func ToOrderResponse(order models.Order) OrderResponse {
	var orderResponse OrderResponse
	orderResponse.OrderID = order.ID
	orderResponse.CustomerName = order.CustomerName
	orderResponse.Status = order.Status
	orderResponse.CreatedAt = order.CreatedAt

	for _, item := range order.OrderItems {
		var itemRequest OrderItemsRequest
		itemRequest.ProductID = item.ProductID
		itemRequest.Quantity = item.Quantity
		orderResponse.Items = append(orderResponse.Items, itemRequest)
	}
	return orderResponse
}
