package dto

import (
	"database/sql"
	"errors"
	"inventory-service/internal/adapter/postgres/dao"
	"net/http"

	"github.com/jackc/pgx/v5"
)

var (
	ErrInvalidFilters      = errors.New("invalid filters")
	ErrUnprocessableEntity = errors.New("unprocessable entity")
	ErrEditConflict        = errors.New("unable to update the record due to an edit conflict")
)

type HTTPError struct {
	Code    int
	Message string
}

var (
	ErrPageNotFoundResponse = &HTTPError{
		Code:    http.StatusNotFound,
		Message: "the requested resource could not be found",
	}
	ErrUnprocessableEntityResponse = &HTTPError{
		Code:    http.StatusUnprocessableEntity,
		Message: "unprocessable entity",
	}
	ErrEditConflictResponse = &HTTPError{
		Code:    http.StatusConflict,
		Message: "unable to update the record due to an edit conflict",
	}
)

func FromError(err error) *HTTPError {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrPageNotFoundResponse
	case errors.Is(err, pgx.ErrNoRows):
		return ErrPageNotFoundResponse
	case errors.Is(err, dao.ErrRecordNotFound):
		return ErrPageNotFoundResponse
	case errors.Is(err, ErrUnprocessableEntity):
		return ErrUnprocessableEntityResponse
	case errors.Is(err, ErrEditConflict):
		return ErrEditConflictResponse
	default:
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "something went wrong",
		}
	}
}