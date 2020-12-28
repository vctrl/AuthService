package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vctrl/authService/middleware"
	"golang.org/x/oauth2"
)

const (
	testSite    = "testSite"
	unknownSite = "unknownSite"
	testURL     = "testURL"
	testCode    = "testCode"

	testAccessToken  = "testAccessToken"
	testTokenType    = "testTokenType"
	testRefreshToken = "testRefreshToken"
)

func TestLoginHappyCase(t *testing.T) {
	configMock := new(OAuth2ConfigMock)
	configMock.On("AuthCodeURL", mock.Anything, mock.Anything).Return(testURL)

	uc := NewAuthUseCaseConfig(testSite, configMock)
	url, _, err := uc.Login(testSite)

	assert.Equal(t, url, testURL)
	assert.Equal(t, err, nil)
}

func TestLoginUnknownSiteReturnError(t *testing.T) {
	configMock := new(OAuth2ConfigMock)
	// mock.Anything because state is random
	configMock.On("AuthCodeURL", mock.Anything, mock.Anything).Return(testURL)

	uc := NewAuthUseCaseConfig(testSite, configMock)
	_, _, err := uc.Login(unknownSite)

	assert.Equal(t, err, &middleware.AppError{Code: 400, Message: "Unknown site"})
}

func TestCallbackUsecase(t *testing.T) {
	configMock := new(OAuth2ConfigMock)
	testToken := &oauth2.Token{
		AccessToken:  testAccessToken,
		TokenType:    testTokenType,
		RefreshToken: testRefreshToken,
		Expiry:       time.Time{},
	}

	expected := &Response{
		AccessToken:  testAccessToken,
		TokenType:    testTokenType,
		RefreshToken: testRefreshToken,
		Expiry:       time.Time{},
	}

	configMock.On("Exchange", oauth2.NoContext, testCode, mock.Anything).Return(testToken, nil)

	uc := NewAuthUseCaseConfig(testSite, configMock)
	r, err := uc.Callback(testSite, testCode)
	assert.Equal(t, r, expected)
	assert.Equal(t, err, nil)
}

func TestCallbackUnknownSiteReturnError(t *testing.T) {
	configMock := new(OAuth2ConfigMock)
	uc := NewAuthUseCaseConfig(testSite, configMock)
	_, err := uc.Callback(unknownSite, testCode)

	assert.Equal(t, err, &middleware.AppError{Code: 400, Message: "Unknown site"})
}

func TestCallbackExchangeTokenReturnError(t *testing.T) {
	configMock := new(OAuth2ConfigMock)
	expectedErr := errors.New("failed")
	configMock.On("Exchange", oauth2.NoContext, testCode, mock.Anything).Return(new(oauth2.Token), expectedErr)
	uc := NewAuthUseCaseConfig(testSite, configMock)
	_, err := uc.Callback(testSite, testCode)
	assert.Equal(t, expectedErr, err)
}

func TestMappingSiteConfigToOauth2Config(t *testing.T) {
	//todo
}
