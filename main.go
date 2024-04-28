package main

import (
	"gantre/handlers"
	"gantre/middlewares"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lpernett/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panicln(err)
		return
	}
}

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templ/**/*")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "Cache-Control", "Refresh"},
		ExposeHeaders:    []string{"Content-Length", "Authorization", "Refresh-Token", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/static", "./static")

	r.POST("/admit", handlers.ClearCookie)

	userGroup := r.Group("/:est_type")
	{
		userGroup.GET("/:est_name", middlewares.RequestParams, middlewares.CheckCookies, handlers.QueueHandler)
		userGroup.GET("/:est_name/ws/:user", middlewares.RequestParams, middlewares.CheckCookies, handlers.WebsocketHandler)
	}

	// qm is short for queue manager
	clientGroup := r.Group("/qm/:est_type")
	{
		clientGroup.GET("/:est_name", middlewares.RequestParams, handlers.FrontDeskHandler)
		clientGroup.POST("/:est_name", middlewares.RequestParams, handlers.UpdateFromFrontDeskHandler)
		r.Run()
	}
}
