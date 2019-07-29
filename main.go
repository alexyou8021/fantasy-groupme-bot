package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

type msg struct {
	attachments map[string]string
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
	router.POST("/", func(c *gin.Context) {
		x, _ := ioutil.ReadAll(c.Request.Body)
		message, _ := json.Marshal(x)
        	log.Printf("%s-----", string(x))
        	log.Println(message)
		c.JSON(http.StatusOK, c)
	})

	router.Run(":" + port)
}
