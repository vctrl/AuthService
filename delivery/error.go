package delivery

import (
	"net/http"

	"github.com/vctrl/authService/middleware"
)

func NewNoSiteError() *middleware.AppError {
	return &middleware.AppError{
		Code:    http.StatusBadRequest,
		Message: "No site parameter is provided",
	}
}

func NewInvalidStateError() *middleware.AppError {
	return &middleware.AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid state",
	}
}

// NewErrNoCookie is the wrapper for ErrNoCookie from http package.
func NewErrNoCookie(err error) *middleware.AppError {
	return &middleware.AppError{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	}
}
