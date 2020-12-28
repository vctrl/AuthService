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
	testState  = "test_state"
	testState2 = "test_state2"

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
	req := createRequest(testState, testState, false)

	uc.On("Callback", testSite, testAuthCode).Return(&usecase.Response{}, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCallbackUsecaseReturnError(t *testing.T) {
	r, uc := createTestContext()
	req := createRequest(testState, testState, true)

	uc.On("Callback", testSite, testAuthCode).Return(&usecase.Response{}, errors.New("failed"))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCallbackHappyCase(t *testing.T) {
	r, uc := createTestContext()
	req := createRequest(testState, testState, true)

	uc.On("Callback", testSite, testAuthCode).Return(&usecase.Response{}, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCallbackInvalidStateReturnError(t *testing.T) {
	r, uc := createTestContext()

	// states doesn't match
	req := createRequest(testState, testState2, true)

	uc.On("Callback", testSite, testAuthCode).Return(&usecase.Response{}, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
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

func createRequest(testState1, testState2 string, withCookie bool) *http.Request {
	reqPath := fmt.Sprintf("/callback?state=%s", testState)
	req, _ := http.NewRequest("GET", reqPath, nil)

	q := req.URL.Query()
	q.Add("state", testState1)
	q.Add("site", testSite)
	q.Add("code", testAuthCode)
	req.URL.RawQuery = q.Encode()

	if withCookie {
		cookie := &http.Cookie{Name: oauthStateHeader, Value: testState2}
		req.AddCookie(cookie)
	}

	return req
}
