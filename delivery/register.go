package delivery

import (
	"net/http"

	"github.com/vctrl/authService/usecase"
)

func RegisterHTTPEndpoints(uc *usecase.OAuthUseCase) {
	h := NewHandler(uc)

	// http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", h.Login)
	http.HandleFunc("/callback", h.Callback)
}
