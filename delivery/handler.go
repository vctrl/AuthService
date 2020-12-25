package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vctrl/authService/usecase"
)

var (
	oauthStateHeader = "oauth_state"
)

type Handler struct {
	usecase *usecase.OAuthUseCase
}

func NewHandler(usecase *usecase.OAuthUseCase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	site := r.FormValue("site")
	if site == "" {
		// todo
		return
	}

	url, state, err := h.usecase.Login(site)
	if err != nil {
		// todo
		return
	}

	fmt.Println(url, state)
	expiration := time.Now().Add(365 * 24 * time.Hour)

	cookie := http.Cookie{Name: oauthStateHeader, Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	// Ensure that there is no request forgery going on, and that the user
	// sending us this connect request is the user that was supposed to.
	state, err := r.Cookie(oauthStateHeader)

	if err != nil {
		// todo write error
		return
	}

	if r.FormValue("state") != state.Value {
		// todo write error
		return
	}

	res, err := h.usecase.Callback(r.FormValue("site"), r.FormValue("code"))
	if err != nil {
		// todo write error
		return
	}

	result, err := json.Marshal(res)
	if err != nil {
		// todo write error
		return
	}

	w.Write(result)
}
