package handlers

import (
	"context"
	"order-service/internal/models"
)

type OrderUsecase interface {
	Create(ctx context.Context, request models.Order) (models.OrderResponce, error)
	Get(ctx context.Context, id int64) (models.Order, error)
	GetList(ctx context.Context) ([]models.Order, error)
	SetStatus(ctx context.Context, request models.UpdateStatus) (models.Order, error)
}
