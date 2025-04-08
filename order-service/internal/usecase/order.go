package usecase

import (
	"context"
	"order-service/internal/models"
)

type Order struct {
	orderRepo OrderRepository
}

func NewOrder(orderRepo OrderRepository) *Order {
	return &Order{
		orderRepo: orderRepo,
	}
}

func (u *Order) Create(ctx context.Context, request models.Order) (models.Order, error) {
	// fmt.Println(request.OrderItems[0].Quantity)
	orderID, err := u.orderRepo.CreateOrder(ctx, request)
	if err != nil {
		return models.Order{}, err
	}
	request.ID = orderID
	return request, nil
}

func (u *Order) Get(ctx context.Context, id int64) (models.Order, error) {

	order, err := u.orderRepo.GetByFilter(ctx, models.OrderFilter{ID: id})
	if err != nil {
		return models.Order{}, err
	}
	return order, nil
}

func (u *Order) GetList(ctx context.Context) ([]models.Order, error) {
	orders, err := u.orderRepo.GetListOrderByFilter(ctx, models.OrderFilter{})
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (u *Order) SetStatus(ctx context.Context, request models.UpdateStatus) (models.Order, error) {
	order, err := u.Get(ctx, request.OrderID)
	if err != nil {
		return models.Order{}, nil
	}
	var updatedOrder models.OrderUpdateData
	updatedOrder.ID = &order.ID
	updatedOrder.Status = &request.Status
	err = u.orderRepo.UpdateOrder(ctx, updatedOrder)
	if err != nil {
		return models.Order{}, err
	}
	order.Status = request.Status

	return order, nil
}
