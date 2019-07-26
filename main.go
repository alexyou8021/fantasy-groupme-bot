package main

import (
	"log"
	"net/http"
	"os"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		fmt.Println(c.Request)
		c.String(http.StatusOK, "success")
	})

	router.Run(":" + port)
}
