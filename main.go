package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

type msg struct {
	text string `json:"text"`
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
	router.POST("/", func(c *gin.Context) {
		var json msg
		if c.BindJSON(&json) == nil {
        		log.Println("-----")
        		log.Println(json)
			c.JSON(http.StatusOK, nil)
		}
	})

	router.Run(":" + port)
}
