package usecase

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	endpoints = map[string]oauth2.Endpoint{"google": google.Endpoint}
)

type Response struct {
	UserInfo     *UserInfo
	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       time.Time
}

type UserInfo struct {
	Id            string `json: id`
	Email         string `json: email`
	VerifiedEmail bool   `json: verified_email`
	Picture       string `json: picture`
}

type SiteConfigs map[string]SiteConfig

type SiteConfig struct {
	ClientIDEnvVar     string
	ClientSecretEnvVar string
	RedirectURL        string
	Scopes             []string
}

func (s *SiteConfig) mapToOauth2Config(site string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv(s.ClientIDEnvVar),
		ClientSecret: os.Getenv(s.ClientSecretEnvVar),
		Endpoint:     endpoints[site],
		RedirectURL:  s.RedirectURL,
		Scopes:       s.Scopes,
	}
}

func NewAuthUseCase() *OAuthUseCase {
	var siteConfigs SiteConfigs

	err := viper.Unmarshal(&siteConfigs)
	configs := make(map[string]*oauth2.Config)
	for s, c := range siteConfigs {
		configs[s] = c.mapToOauth2Config(s)
	}

	if err != nil {
		fmt.Println("Error reading configs: ", err)
	}

	return &OAuthUseCase{
		Configs: configs,
	}
}

type OAuthUseCase struct {
	Configs map[string]*oauth2.Config
}

func (auth *OAuthUseCase) getSiteConfig(site string) (*oauth2.Config, error) {
	c, ok := auth.Configs[site]

	if !ok {
		return nil, ErrUnknownSite
	}

	return c, nil
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

	// todo move to config
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	ui := &UserInfo{}
	err = json.Unmarshal(content, ui)
	if err != nil {
		return nil, err
	}

	// todo return just user info
	return &Response{
		UserInfo:     ui,
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry}, nil
}
