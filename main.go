package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)
type Info struct {
	attachments []map[string]string
	avatar_url string
	created_at int
	group_id string
	id string
	name string
	sender_id string
	sender_type string
	source_guid string
	system bool
	text string
	user_id string
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
	router.POST("/" , func(c *gin.Context) {
		var data Info
        	c.BindJSON(&data)
		log.Println(data.name)
	})

	router.Run(":" + port)
}
