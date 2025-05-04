package middleware

import (
	"GO/internal/config"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

var store *sessions.CookieStore

// INFO プログラム起動時に実行
func init() {
	cfg := config.GetConfig()
	store = sessions.NewCookieStore([]byte(cfg.OAuth.SecretKey))
}

func RequireLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		c.Response()
		session, _ := store.Get(req, "session-name")
		user := session.Values["user"]
		if user == nil {
			return c.Redirect(http.StatusFound, "/login")
		}
		return next(c)
	}
}

func SaveUserSession(c echo.Context, username string) {
	req := c.Request()
	res := c.Response()
	session, _ := store.Get(req, "session-name")
	session.Values["user"] = username
	session.Save(req, res)
}

func GetSession(c echo.Context) (*sessions.Session, error) {
	req := c.Request()
	return store.Get(req, "session-name")
}
