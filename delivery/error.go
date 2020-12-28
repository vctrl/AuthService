package delivery

import (
	"net/http"

	"github.com/vctrl/authService/middleware"
)

func NewUnknownSiteError() *middleware.AppError {
	return &middleware.AppError{
		Code:    http.StatusBadRequest,
		Message: "Unknown site",
	}
}

func NewNoSiteError() *middleware.AppError {
	return &middleware.AppError{
		Code:    http.StatusBadRequest,
		Message: "No site is provided",
	}
}

func NewInvalidStateError() *middleware.AppError {
	return &middleware.AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid state",
	}
}
