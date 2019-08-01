package main

import (
	"log"
	"net/http"
	"os"
        "io/ioutil"
        "encoding/json"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

type test struct {
    Response map[string]interface{} `json:"response"`
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
                url := "https://api.groupme.com/v3/groups/18129715?token="
                url = url + os.Getenv("token")
                resp, _ := http.Get(url)

                defer resp.Body.Close()
                bodyBytes, _ := ioutil.ReadAll(resp.Body)
                var test1 test
                json.Unmarshal(bodyBytes, &test1)
                members := test1.Response["members"][0]["nickname"]
                log.Println(test1)
                log.Println(members)
                
		c.String(http.StatusOK, "success")
	})
	router.POST("/", msgHandler())
	router.POST("/reminders", reminderHandler())

	router.Run(":" + port)
}
