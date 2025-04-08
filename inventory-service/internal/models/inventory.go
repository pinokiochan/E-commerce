package models

import "time"

type Inventory struct {
	ID          int64
	Name        string
	Description string
	Price       float64
	Available   int64
	CreatedAt   time.Time
	Version     int32
	IsDeleted   bool
}

type UpdateInventoryData struct {
	ID          *int64
	Name        *string
	Description *string
	Price       *float64
	Available   *int64
	CreatedAt   *time.Time
	Version     *int32
	IsDeleted   *bool
}
