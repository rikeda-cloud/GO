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

func initCarData() {
	carDataDB.InsertCarData("1.png", 0.1, -30)
	carDataDB.InsertCarData("2.png", 0.2, -20)
	carDataDB.InsertCarData("3.png", 0.3, -10)
	carDataDB.InsertCarData("4.png", 0.4, 0)
	carDataDB.InsertPredictedCarData("5.png", 0.5, 10)
	carDataDB.InsertPredictedCarData("6.png", 0.6, 20)
	carDataDB.InsertPredictedCarData("7.png", 0.7, 30)
	carDataDB.InsertPredictedCarData("8.png", 0.8, 40)
	carDataDB.InsertPredictedCarData("9.png", 0.9, 50)
}

func main() {
	cfg := config.GetConfig()
	if err := carDataDB.CreateCarDataTableIf(); err != nil {
		log.Fatal(err)
	}
	initCarData()

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
	e.GET("/logout", func(c echo.Context) error {
		session, _ := middleware.GetSession(c)
		session.Options.MaxAge = -1
		session.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusFound, "/login")
	})
	e.Logger.Fatal(e.Start(cfg.App.Annotation.Port))
}
