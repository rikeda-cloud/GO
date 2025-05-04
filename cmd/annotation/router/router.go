package router

import (
	"GO/cmd/annotation/handlers"
	"GO/internal/config"
	"GO/internal/middleware"
	"GO/internal/oauth"
	"net/http"
	"path"
	"slices"

	"github.com/labstack/echo/v4"
)

func SetupRouter() *echo.Echo {
	cfg := config.GetConfig()

	e := echo.New()

	// INFO 認証に関連したエンドポイント
	e.GET("/login", handleLogin)
	e.GET("/callback", handleCallback)
	e.GET("/api/user", handleAuthUser)

	// INFO authGroupに所属するエンドポイントは認証済み出ないとアクセス不可
	authGroup := e.Group("")
	authGroup.Use(middleware.RequireLogin)

	authGroup.GET("/", handleIndexHtml)
	authGroup.GET("/*", echo.WrapHandler(http.StripPrefix("/", http.FileServer(http.Dir(cfg.App.Annotation.StaticDir)))))
	authGroup.Static("/images/", cfg.Image.DirPath)
	authGroup.Static("/predict-images/", cfg.Image.PredictDirPath)

	authGroup.GET("/ws", handlers.NewAnnotationHandler().HandleAnnotation)
	authGroup.GET("/ws/remain-count", handlers.NewRemainImageCountHandler().HandleRemainImageCount)
	authGroup.GET("/ws/predict-remain-count", handlers.NewPredictedRemainImageCountHandler().HandlePredictedRemainImageCount)
	authGroup.GET("/ws/check", handlers.NewAnnotatedDataCheckHandler().HandleAnnotatedDataCheck)
	authGroup.GET("/ws/ai", handlers.NewPredictedDataHandler().HandlePredictedData)

	return e
}

func handleLogin(c echo.Context) error {
	url := oauth.GetLoginURL("state")
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func handleCallback(c echo.Context) error {
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

func handleIndexHtml(c echo.Context) error {
	cfg := config.GetConfig()
	return c.File(path.Join(cfg.App.Annotation.StaticDir, "index.html"))
}

func handleAuthUser(c echo.Context) error {
	session, _ := middleware.GetSession(c)
	user := session.Values["user"]
	if user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthenticated"})
	}
	return c.JSON(http.StatusOK, map[string]string{"user": user.(string)})
}
