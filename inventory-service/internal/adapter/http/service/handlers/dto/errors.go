package dto

import (
	"database/sql"
	"errors"
	"inventory-service/internal/adapter/postgres/dao"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type HTTPError struct {
	Code    int
	Message string
}

var (
	ErrPageNotFound = &HTTPError{
		Code:    http.StatusNotFound,
		Message: "The requested page could not be found",
	}
)

func FromError(err error) *HTTPError {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrPageNotFound
	case errors.Is(err, pgx.ErrNoRows):
		return ErrPageNotFound
	case errors.Is(err, dao.ErrRecordNotFound):
		return ErrPageNotFound
	default:
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "something went wrong",
		}

	}

}
