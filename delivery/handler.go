package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (h *Handler) Login(c *gin.Context) {
	site := c.Request.FormValue("site")
	if site == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	url, state, err := h.usecase.Login(site)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie(oauthStateHeader, state, 60*60*24, "/", "localhost", false, true)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) Callback(c *gin.Context) {
	// Ensure that there is no request forgery going on, and that the user
	// sending us this connect request is the user that was supposed to.
	state, err := c.Cookie(oauthStateHeader)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if c.Request.FormValue("state") != state {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res, err := h.usecase.Callback(c.Request.FormValue("site"), c.Request.FormValue("code"))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(res)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Writer.Write(result)
}
