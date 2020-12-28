package usecase

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
