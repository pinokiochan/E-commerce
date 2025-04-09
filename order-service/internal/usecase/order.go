package usecase

import (
	"context"
	"fmt"
	"order-service/internal/models"
)

type Order struct {
	orderRepo        OrderRepository
	inventoryService InventoryService
}

func NewOrder(orderRepo OrderRepository, inventoryService InventoryService) *Order {
	return &Order{
		orderRepo:        orderRepo,
		inventoryService: inventoryService,
	}
}

func (u *Order) Create(ctx context.Context, request models.Order) (models.OrderResponce, error) {
	// Inserting order to database
	orderID, err := u.orderRepo.Create(ctx, request)
	if err != nil {
		return models.OrderResponce{}, err
	}

	// Metadata of items
	var orderItemResponces []models.OrderItemResponce
	var totalPrice int64
	for _, item := range request.OrderItems {
		var orderItemResp models.OrderItemResponce
		orderItemResp.ProductID = item.ProductID

		// Getting inventory from inventory service
		inventoryItem, err := u.inventoryService.GetById(item.ProductID)
		if err != nil {
			orderItemResp.Status = "rejected"
			orderItemResp.Reason = err.Error()
			orderItemResponces = append(orderItemResponces, orderItemResp)
			continue
		}

		price := inventoryItem.Price * item.Quantity
		newAvailability := inventoryItem.Available - item.Quantity

		if newAvailability < 0 {
			orderItemResp.Status = "rejected"
			orderItemResp.Reason = "insufficient_inventory"
			orderItemResponces = append(orderItemResponces, orderItemResp)
			continue
		}

		// Trying to set new availability
		err = u.inventoryService.Substruct(item.ProductID, newAvailability, inventoryItem.Version)
		if err != nil {
			orderItemResp.Status = "rejected"
			orderItemResp.Reason = err.Error()
			orderItemResponces = append(orderItemResponces, orderItemResp)
			continue
		}

		orderItemResp.Name = inventoryItem.Name
		orderItemResp.Price = price
		orderItemResp.Status = "accepted"

		totalPrice += price

		orderItemResponces = append(orderItemResponces, orderItemResp)
	}

	fmt.Println(orderItemResponces)

	responce := models.OrderResponce{
		OrderID:      orderID,
		CustomerName: request.CustomerName,
		Items:        orderItemResponces,
		Total:        totalPrice,
	}

	return responce, nil
}

func (u *Order) GetList(ctx context.Context) ([]models.Order, error) {
	orders, err := u.orderRepo.GetListWithFilter(ctx, models.OrderFilter{})
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (u *Order) Get(ctx context.Context, id int64) (models.Order, error) {
	order, err := u.orderRepo.GetWithFilter(ctx, models.OrderFilter{ID: id})
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (u *Order) SetStatus(ctx context.Context, req models.UpdateStatus) (models.Order, error) {
	order, err := u.Get(ctx, req.OrderID)
	if err != nil {
		return models.Order{}, err
	}

	var updatedOrder models.OrderUpdateData
	updatedOrder.ID = &order.ID
	updatedOrder.Status = &req.Status

	err = u.orderRepo.Update(ctx, updatedOrder)

	if err != nil {
		return models.Order{}, err
	}

	order.Status = req.Status
	return order, nil
}
