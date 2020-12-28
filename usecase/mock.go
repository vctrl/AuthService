package usecase

import (
	"github.com/stretchr/testify/mock"
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
