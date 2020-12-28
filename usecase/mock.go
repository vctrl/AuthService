package usecase

import (
	"context"

	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

type OAuthUseCaseMock struct {
	mock.Mock
}

func (m *OAuthUseCaseMock) Login(site string) (string, string, error) {
	args := m.Called(site)
	return args.Get(0).(string), args.Get(1).(string), args.Error(2)
}

func (m *OAuthUseCaseMock) Callback(site, authCode string) (*Response, error) {
	args := m.Called(site, authCode)
	return args.Get(0).(*Response), args.Error(1)
}

type OAuth2ConfigMock struct {
	mock.Mock
}

func (m *OAuth2ConfigMock) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	args := m.Called(state, opts)
	return args.Get(0).(string)
}

func (m *OAuth2ConfigMock) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	args := m.Called(ctx, code, opts)
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (m *OAuth2ConfigMock) MapToOauth2Config(site string) *oauth2.Config {
	args := m.Called(site)
	return args.Get(0).(*oauth2.Config)
}
