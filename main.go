package main

import (
	"gantre/handlers"
	"gantre/middlewares"
	"log"

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

	userGroup := r.Group("/:est_type")
	{
		userGroup.GET("/:est_name", middlewares.RequestParams, middlewares.CheckCookies, handlers.QueueHandler)
	}

	r.POST("/test", handlers.InitQueue)

	r.Run()
}
