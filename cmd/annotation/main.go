package main

import (
	"GO/cmd/annotation/handlers"
	"GO/internal/config"
	"GO/internal/db"
	"GO/internal/middleware"
	"GO/internal/oauth"
	"log"
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slices"
)

func main() {
	cfg := config.GetConfig()
	if err := carDataDB.CreateCarDataTableIf(); err != nil {
		log.Fatal(err)
	}
	carDataDB.InitTmpCarData()

	e := echo.New()

	e.GET("/login", func(c echo.Context) error {
		url := oauth.GetLoginURL("state")
		return c.Redirect(http.StatusTemporaryRedirect, url)
	})

	e.GET("/callback", func(c echo.Context) error {
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
	})

	authGroup := e.Group("")
	authGroup.Use(middleware.RequireLogin)

	authGroup.GET("/", func(c echo.Context) error {
		return c.File(path.Join(cfg.App.Annotation.StaticDir, "index.html"))
	})
	authGroup.GET("/*", echo.WrapHandler(http.StripPrefix("/", http.FileServer(http.Dir(cfg.App.Annotation.StaticDir)))))
	authGroup.Static("/images/", cfg.Image.DirPath)
	authGroup.Static("/predict-images/", cfg.Image.PredictDirPath)
	authGroup.GET("/ws", handlers.NewAnnotationHandler().HandleAnnotation)
	authGroup.GET("/ws/remain-count", handlers.NewRemainImageCountHandler().HandleRemainImageCount)
	authGroup.GET("/ws/predict-remain-count", handlers.NewPredictedRemainImageCountHandler().HandlePredictedRemainImageCount)
	authGroup.GET("/ws/check", handlers.NewAnnotatedDataCheckHandler().HandleAnnotatedDataCheck)
	authGroup.GET("/ws/ai", handlers.NewPredictedDataHandler().HandlePredictedData)
	e.GET("/api/user", func(c echo.Context) error {
		session, _ := middleware.GetSession(c)
		user := session.Values["user"]
		if user == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthenticated"})
		}
		return c.JSON(http.StatusOK, map[string]string{"user": user.(string)})
	})
	e.Logger.Fatal(e.Start(cfg.App.Annotation.Port))
}
