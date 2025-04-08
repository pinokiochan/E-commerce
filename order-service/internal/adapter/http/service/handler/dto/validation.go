package dto

import (
	"fmt"
	"order-service/internal/adapter/postgres/dao"
	"order-service/internal/models"
	"order-service/pkg/validator"
	"strings"
)

func ValidateOrder(e *validator.Validator, order models.Order) {
	e.Check(order.CustomerName != "", "customer_name", "must be provided")
	e.Check(len(order.CustomerName) < 50, "customer_name", "must be more than 50 bytes long")

	for _, item := range order.OrderItems {
		e.Check(item.ProductID > 0, "items_product_id", "must be greater than zero")
		e.Check(item.Quantity > 0, "items_quantity", "must be greater than zero")
		e.Check(item.Quantity <= 100, "items_quantity", "item quantity cannot be greater than 100")
	}

}

func ValidateSetOrderStatusRequest(v *validator.Validator, req OrderStatusRequest) {
	safeList := []string{dao.OrderStatusCanceled, dao.OrderStatusCompleted, dao.OrderStatusPending}
	v.Check(validator.PermittedValue(req.Status, safeList...), "status", fmt.Sprintf("invalid status. Available: %v", strings.Join(safeList, ", ")))
}