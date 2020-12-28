package delivery

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vctrl/authService/middleware"
	"github.com/vctrl/authService/usecase"
)

const (
	testState    = "test_state"
	testSite     = "test_site"
	testAuthCode = "test_auth_code"
)

func TestLoginWithNoSiteParamReturnsBadRequest(t *testing.T) {
	r, _ := createTestContext()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginReturnError(t *testing.T) {
	r, uc := createTestContext()

	uc.On("Login", "site").Return("", "", errors.New("failed"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login?site=site", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestLoginHappyCase(t *testing.T) {
	r, uc := createTestContext()

	uc.On("Login", "site").Return("test_url", "test_state", nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login?site=site", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
}

func TestCallbackNoCookieReturnError(t *testing.T) {
	r, uc := createTestContext()
	reqPath := fmt.Sprintf("/callback?state=%s", testState)
	req, _ := http.NewRequest("GET", reqPath, nil)

	q := req.URL.Query()
	q.Add("state", testState)
	q.Add("site", testSite)
	q.Add("code", testAuthCode)
	req.URL.RawQuery = q.Encode()

	uc.On("Callback", testSite, testAuthCode).Return(&usecase.Response{}, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	fmt.Println(w)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCallbackUsecaseReturnError(t *testing.T) {
	r, uc := createTestContext()
	reqPath := fmt.Sprintf("/callback?state=%s", testState)
	req, _ := http.NewRequest("GET", reqPath, nil)
	cookie := &http.Cookie{Name: oauthStateHeader, Value: testState}
	req.AddCookie(cookie)

	q := req.URL.Query()
	q.Add("state", testState)
	q.Add("site", testSite)
	q.Add("code", testAuthCode)
	req.URL.RawQuery = q.Encode()

	uc.On("Callback", testSite, testAuthCode).Return(&usecase.Response{}, errors.New("failed"))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCallbackHappyCase(t *testing.T) {
	r, uc := createTestContext()
	reqPath := fmt.Sprintf("/callback?state=%s", testState)
	req, _ := http.NewRequest("GET", reqPath, nil)
	cookie := &http.Cookie{Name: oauthStateHeader, Value: testState}
	req.AddCookie(cookie)

	q := req.URL.Query()
	q.Add("state", testState)
	q.Add("site", testSite)
	q.Add("code", testAuthCode)
	req.URL.RawQuery = q.Encode()

	uc.On("Callback", testSite, testAuthCode).Return(&usecase.Response{}, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func createTestContext() (*gin.Engine, *usecase.OAuthUseCaseMock) {
	r := gin.Default()
	r.Use(
		middleware.JSONAppErrorReporter(),
	)

	uc := new(usecase.OAuthUseCaseMock)
	RegisterHTTPEndpoints(r, uc)
	return r, uc
}
