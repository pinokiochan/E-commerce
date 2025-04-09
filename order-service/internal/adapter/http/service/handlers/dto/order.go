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

type OrderCreateResponceRequest struct {
	OrderID      int64  `json:"order_id"`
	CustomerName string `json:"customer_name"`
}

type OrderCreateResponceRequestV2 struct {
	OrderID      int64                               `json:"order_id"`
	CustomerName string                              `json:"customer_name"`
	Items        []OrderItemsCreateResponceRequestV2 `json:"items"`
	Total        int64                               `json:"total"`
}

type OrderItemsCreateResponceRequestV2 struct {
	ProductID int64  `json:"product_id"`
	Name      string `json:"name"`
	Price     int64  `json:"price,omitempty"`  // Total price
	Status    string `json:"status,omitempty"` // accepted, rejected
	Reason    string `json:"reason,omitempty"` // if rejected
}

type OrderResponce struct {
	OrderID      int64               `json:"order_id"`
	CustomerName string              `json:"customer_name"`
	Items        []OrderItemsRequest `json:"items"`
	Status       string              `json:"status"`
	CreatedAt    time.Time           `json:"created_at"`
}

type OrderSetStatusRequest struct {
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

func ToOrderCreateResponse(order models.OrderResponce) OrderCreateResponceRequestV2 {
	var itemsInfo []OrderItemsCreateResponceRequestV2

	for _, v := range order.Items {
		itemsInfo = append(itemsInfo, OrderItemsCreateResponceRequestV2{
			ProductID: v.ProductID,
			Name:      v.Name,
			Price:     v.Price,
			Status:    v.Status,
			Reason:    v.Reason,
		})
	}

	return OrderCreateResponceRequestV2{
		OrderID:      order.OrderID,
		CustomerName: order.CustomerName,
		Items:        itemsInfo,
		Total:        order.Total,
	}
}

func ToOrderListResponce(orders []models.Order) []OrderResponce {
	resp := []OrderResponce{}

	for _, order := range orders {
		orderResponce := ToOrderResponce(order)
		resp = append(resp, orderResponce)
	}

	return resp
}

func ToOrderResponce(order models.Order) OrderResponce {
	var orderResponce OrderResponce

	orderResponce.OrderID = order.ID
	orderResponce.CustomerName = order.CustomerName
	orderResponce.Status = order.Status
	orderResponce.CreatedAt = order.Created_at

	for _, item := range order.OrderItems {
		var itemRequest OrderItemsRequest
		itemRequest.ProductID = item.ProductID
		itemRequest.Quantity = item.Quantity
		orderResponce.Items = append(orderResponce.Items, itemRequest)
	}

	return orderResponce
}
