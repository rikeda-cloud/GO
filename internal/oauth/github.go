package oauth

import (
	"GO/internal/config"
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var cfg = config.GetConfig()

var oauthConf = &oauth2.Config{
	ClientID:     cfg.OAuth.ClientID,
	ClientSecret: cfg.OAuth.ClientSecret,
	Scopes:       []string{"user:email"},
	Endpoint:     github.Endpoint,
	RedirectURL:  cfg.OAuth.CallbackURL,
}

func GetLoginURL(state string) string {
	return oauthConf.AuthCodeURL(state)
}

func GetGitHubUserName(code string) (string, error) {
	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		return "", err
	}
	client := oauthConf.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var user struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", err
	}
	if user.Login == "" {
		return "", errors.New("empty GitHub login")
	}
	return user.Login, nil
}
