package dao

import (
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")

	SafeSortList = []string{"id", "name", "price", "available", "-id", "-name", "-price", "-available"}
)

type Product struct {
	ID          int64
	CreatedAt   time.Time
	Name        string
	Description string
	Price       float64
	Available   *int

	IsDeleted bool
	Version   int
}
