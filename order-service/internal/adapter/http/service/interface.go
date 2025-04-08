package service

import "order-service/internal/adapter/http/service/handler"

type OrderUsecase interface {
	handler.OrderUsecase
}