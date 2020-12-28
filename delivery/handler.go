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
	usecase usecase.UseCase
}

func NewHandler(usecase usecase.UseCase) *Handler {
	return &Handler{usecase: usecase}
}

// Get redirection URL to the authorization service
func (h *Handler) Login(c *gin.Context) {
	site := c.Request.FormValue("site")
	if site == "" {
		c.Error(NewNoSiteError())
		return
	}

	url, state, err := h.usecase.Login(site)
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookie(oauthStateHeader, state, 60*60*24, "/", "localhost", false, true)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Callback URL to which the authorization service redirects.
func (h *Handler) Callback(c *gin.Context) {
	// Ensure that there is no request forgery going on, and that the user
	// sending us this connect request is the user that was supposed to.
	state, err := c.Cookie(oauthStateHeader)

	if err != nil {
		c.Error(NewErrNoCookie(err))
		return
	}

	if c.Request.FormValue("state") != state {
		c.Error(NewInvalidStateError())
		return
	}

	res, err := h.usecase.Callback(c.Request.FormValue("site"), c.Request.FormValue("code"))
	if err != nil {
		c.Error(err)
		return
	}

	result, err := json.Marshal(res)
	if err != nil {
		c.Error(err)
		return
	}

	c.Writer.Write(result)
}
