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
                url := "https://api.groupme.com/v3/groups/18129715?token="
                url = url + os.Getenv("token")
                resp, _ := http.Get(url)

                defer resp.Body.Close()
                bodyBytes, _ := ioutil.ReadAll(resp.Body)
                bodyString := string(bodyBytes)
                log.Println(bodyString)
                
		c.String(http.StatusOK, "success")
	})
	router.POST("/", msgHandler())
	router.POST("/reminders", reminderHandler())

	router.Run(":" + port)
}
