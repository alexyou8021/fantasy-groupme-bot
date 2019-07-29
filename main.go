package main

import (
	"log"
	"net/http"
	"os"
	"io/ioutil"

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
		x, _ := ioutil.ReadAll(c.Request.Body)
        	log.Printf("%s", string(x))
		c.String(http.StatusOK, "success")
	})
	router.POST("/", func(c *gin.Context) {
		x, _ := ioutil.ReadAll(c.Request.Body)
        	log.Printf("%s", string(x))
		c.JSON(http.StatusOK, c)
	})

	router.Run(":" + port)
}
