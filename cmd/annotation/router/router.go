package router

import (
	"GO/cmd/annotation/handlers"
	"GO/internal/config"
	"GO/internal/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupRouter() *echo.Echo {
	cfg := config.GetConfig()

	e := echo.New()

	// INFO 認証に関連したエンドポイント
	e.GET("/login", handlers.HandleLogin)
	e.GET("/callback", handlers.HandleCallback)
	e.GET("/api/user", handlers.HandleAuthUser)

	// INFO authGroupに所属するエンドポイントは認証済み出ないとアクセス不可
	authGroup := e.Group("")
	authGroup.Use(middleware.RequireLogin)

	authGroup.GET("/", handlers.HandleIndexHtml)
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
