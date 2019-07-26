package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

type CreateParams struct {
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})
	router.POST("/login/do", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/welcome")
	})

	router.GET("/welcome", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome")
	})

	router.Run(":" + port)
}
