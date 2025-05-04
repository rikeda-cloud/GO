package handlers

import (
	"GO/internal/config"
	"GO/internal/middleware"
	"GO/internal/oauth"
	"net/http"
	"path"
	"slices"

	"github.com/labstack/echo/v4"
)

func HandleLogin(c echo.Context) error {
	url := oauth.GetLoginURL("state")
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleCallback(c echo.Context) error {
	cfg := config.GetConfig()
	code := c.QueryParam("code")
	user, err := oauth.GetGitHubUserName(code)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Authentication failed")
	}
	if !slices.Contains(cfg.OAuth.AllowedUsers, user) {
		return c.String(http.StatusForbidden, "Access denied")
	}
	middleware.SaveUserSession(c, user)
	return c.Redirect(http.StatusFound, "/")
}

func HandleIndexHtml(c echo.Context) error {
	cfg := config.GetConfig()
	return c.File(path.Join(cfg.App.Annotation.StaticDir, "index.html"))
}

func HandleAuthUser(c echo.Context) error {
	session, _ := middleware.GetSession(c)
	user := session.Values["user"]
	if user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthenticated"})
	}
	return c.JSON(http.StatusOK, map[string]string{"user": user.(string)})
}
