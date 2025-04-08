package dto

import (
	"inventory-service/internal/models"
	"inventory-service/pkg/validator"
)

func ValidateInventory(e *validator.Validator, inv models.Inventory) {
	e.Check(len(inv.Name) != 0, "name", "must be greater than 0")
	e.Check(len(inv.Name) < 50, "name", "must be not greater than 50")
	e.Check(len(inv.Description) != 0, "description", "must be provided")
	e.Check(inv.Price > 0, "price", "must be greater than 0")

}
