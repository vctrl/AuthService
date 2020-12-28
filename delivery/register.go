package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/vctrl/authService/usecase"
)

func RegisterHTTPEndpoints(router *gin.Engine, uc usecase.UseCase) {
	h := NewHandler(uc)

	authEndpoints := router.Group("/")
	{
		authEndpoints.GET("/login", h.Login)
		authEndpoints.GET("/callback", h.Callback)
	}
}
