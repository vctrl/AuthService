package usecase

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	endpoints = map[string]oauth2.Endpoint{"google": google.Endpoint}
)

type UseCase interface {
	Login(site string) (string, string, error)
	Callback(site, authCode string) (*Response, error)
}

type UseCaseConfig interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
}

type Response struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       time.Time
}

type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

type SiteConfigs map[string]SiteConfig

type SiteConfig struct {
	ClientIDEnvVar     string
	ClientSecretEnvVar string
	RedirectURL        string
	Scopes             []string
}

type Config struct {
	Sites SiteConfigs
	Port  int
}

func (s *SiteConfig) MapToOauth2Config(site string, getenv func(string) string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     getenv(s.ClientIDEnvVar),
		ClientSecret: getenv(s.ClientSecretEnvVar),
		Endpoint:     endpoints[site],
		RedirectURL:  s.RedirectURL,
		Scopes:       s.Scopes,
	}
}

func NewAuthUseCase(configs map[string]UseCaseConfig) *OAuthUseCase {
	return &OAuthUseCase{
		Configs: configs,
	}
}

func NewAuthUseCaseConfig(site string, config UseCaseConfig) *OAuthUseCase {
	return &OAuthUseCase{
		Configs: map[string]UseCaseConfig{site: config},
	}
}

type OAuthUseCase struct {
	Configs map[string]UseCaseConfig
}

func (auth *OAuthUseCase) Login(site string) (string, string, error) {
	stateBigInt, err := rand.Int(rand.Reader, big.NewInt(80))
	if err != nil {
		return "", "", err
	}

	stateStr := stateBigInt.String()
	authConfig, err := auth.getSiteConfig(site)
	if err != nil {
		return "", "", err
	}

	url := authConfig.AuthCodeURL(stateStr)
	return url, stateStr, nil
}

func (auth *OAuthUseCase) Callback(site, authCode string) (*Response, error) {
	authConfig, err := auth.getSiteConfig(site)
	if err != nil {
		return nil, err
	}

	token, err := authConfig.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		return nil, err
	}

	return &Response{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry}, nil
}

func (auth *OAuthUseCase) getSiteConfig(site string) (UseCaseConfig, error) {
	c, ok := auth.Configs[site]

	if !ok {
		return nil, NewUnknownSiteError()
	}

	return c, nil
}
