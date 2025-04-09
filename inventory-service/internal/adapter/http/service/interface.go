package service

import "inventory-service/internal/adapter/http/service/handlers"

type InventoryUsecase interface {
	handlers.InventoryUsecase
}
