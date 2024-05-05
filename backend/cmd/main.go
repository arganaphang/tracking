package main

import (
	"application/internals/handler"
	"application/internals/repository"
	"application/internals/service"
	"application/pkg/sse"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	gin.Default().SetTrustedProxies(nil)
	app := gin.New()

	db, err := sqlx.Open("sqlite3", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln("failed to connect to database")
	}

	if db.Ping() != nil {
		log.Fatalln("failed to ping database")
	}

	repositories := repository.Repositories{
		Order: repository.NewOrder(db),
		User:  repository.NewUser(db),
	}
	services := service.Services{
		Order: service.NewOrder(repositories),
		User:  service.NewUser(repositories),
	}

	authMiddleware, err := handler.JWTMiddleware(repositories.User)
	if err != nil {
		log.Fatalln("error to create jwt middleware")
	}
	if err := authMiddleware.MiddlewareInit(); err != nil {
		log.Fatalln("error to initialize jwt middleware")
	}

	sentEvent := sse.NewSSEConn()

	handler.NewUser(app, services, authMiddleware.MiddlewareFunc())
	handler.NewOrder(app, services, sentEvent, authMiddleware.MiddlewareFunc())

	app.GET("/healthz", health)
	app.POST("/login", authMiddleware.LoginHandler)
	app.GET("/auth/refresh-txoken", authMiddleware.RefreshHandler)

	app.Run("0.0.0.0:8000")
}

func health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]any{
		"message": "OK",
	})
}
