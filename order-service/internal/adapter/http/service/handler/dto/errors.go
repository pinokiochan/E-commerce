package dto

import (
	"database/sql"
	"net/http"
	"errors"
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
	default:
		return &HTTPError{
			Code: http.StatusInternalServerError,
			Message: "something went wrong",
		}
	}
}
