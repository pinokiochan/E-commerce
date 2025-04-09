package models

import "time"

type Order struct {
	ID           int64
	CustomerName string
	OrderItems   []OrderItem
	Status       string
	Created_at   time.Time

	IsDeleted bool
}

type OrderItem struct {
	OrderID   int64
	ProductID int64
	Quantity  int64
}

type OrderUpdateData struct {
	ID           *int64
	CustomerName *string
	OrderItems   *[]OrderItem
	Status       *string
	Created_at   *time.Time

	IsDeleted *bool
}

type OrderItemUpdatedData struct {
	OrderID   *int64
	ProductID *int64
	Quantity  *int64
}
type UpdateStatus struct {
	OrderID int64
	Status  string
}

type OrderFilter struct {
	ID int64
}

// OrderInfo
type OrderInfo struct {
	Order         Order
	OrderResponce OrderResponce
}
