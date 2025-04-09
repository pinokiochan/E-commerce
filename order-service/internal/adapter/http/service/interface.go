package service

import "order-service/internal/adapter/http/service/handlers"

type OrderUsecase interface {
	handlers.OrderUsecase
}
