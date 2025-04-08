package dao

import "time"

type Order struct {
	ID           int64
	CustomerName string
	Status       string
	Created_At   time.Time
	isDeleted    bool
}

type OrderItem struct {
	OrderID   int64
	ProductID int64
	Quantity int64
}
