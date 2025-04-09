package models

type Inventory struct {
	ID          int64
	Name        string
	Description string
	Price       int64
	Available   int64
	CreatedAt   string
	Version     int32
}

type OrderResponce struct {
	OrderID      int64
	CustomerName string
	Items        []OrderItemResponce
	Total        int64
}

type OrderItemResponce struct {
	ProductID int64
	Name      string
	Price     int64
	Status    string
	Reason    string
}
