package models

import (
	"inventory-service/pkg/validator"
	"slices"
	"strings"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func ValidateFilters(e *validator.Validator, f Filters) {
	e.Check(f.Page > 0, "page", "must be greater than zero")
	e.Check(f.PageSize <= 10_000_000, "page", "must be maximum of 10 million")
	e.Check(f.PageSize > 0, "page", "must be greater than 0")
	e.Check(f.PageSize <= 100, "page_size", "must be maximum of 100")

	e.Check(validator.PermittedValue(f.Sort, f.SortSafelist...), "sort", "invalid sort value")

}


func (f Filters) SortColumn() string {
	if slices.Contains(f.SortSafelist, f.Sort) {
		return strings.TrimPrefix(f.Sort, "-")
	}
	panic("unsafe sort parameter: " + f.Sort)
}
// Return the sort direction ("ASC" or "DESC") depending on the prefix character of the
// Sort field.
func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) Limit() int {
	return f.PageSize
}

func (f Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}