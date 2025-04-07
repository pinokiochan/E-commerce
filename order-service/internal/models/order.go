package models

import "time"

// main Model
type Order struct {
	ID           int64
	CustomerName string
	OrderItems   []OrderItem //slice of items (one to many)
	Status       string      // status like created, shipped, delivered
	CreatedAt    time.Time
	IsDeleted    bool // soft delete
}

// sub Model
// each position point which (ProductID) and how many (Quantity)
type OrderItem struct { // one order contain several items
	OrderID   int64
	ProductID int64
	Quantity  int64
}

type OrderUpdateData struct { //  * Allows to update only those fields that were passed, ignoring the rest
	ID           *int64
	CustomerName *string
	OrderItems   *[]OrderItem
	Status       *string
	CreatedAt    *time.Time
	IsDeleted    *bool
}

type OrderItemUpdateData struct { // the same idea, update if field isnt nil
	OrderID   *int64
	ProductID *int64
	Quantity  *int64
}
type UpdateStatus struct { // simple dto request to change status
	OrderID int64
	Status  string
}

type OrderFilter struct { // for seacrh so far only ID, can be expanded
	ID int64
}
