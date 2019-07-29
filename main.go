package main

import (
	"log"
	"net/http"
	"encoding/json"
	"bytes"
	"os"
	"math/rand"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

type msg struct {
	Text string `json:"text"`
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
		var botResponse msg
		if c.BindJSON(&botResponse) == nil {
        		log.Println(botResponse.Text)

                        if botResponse.Text == "!help" {
				message := map[string]interface{}{
					"bot_id": os.Getenv("botid"),
					"text": "I am your chat bot. Type `!coin` to flip a coin.",
				}

				bytesRepresentation, err := json.Marshal(message)
				if err != nil {
					log.Fatalln(err)
				}

				http.Post("https://api.groupme.com/v3/bots/post", "application/json", bytes.NewBuffer(bytesRepresentation))
			}
                        if botResponse.Text == "!coin" {
				result := "Your coin landed on HEADS."
				if rand.Intn(2) == 1 {
					result = "Your coin landed on TAILS."
				}
				
				message := map[string]interface{}{
					"bot_id": os.Getenv("botid"),
					"text": result,
				}

				bytesRepresentation, err := json.Marshal(message)
				if err != nil {
					log.Fatalln(err)
				}

				http.Post("https://api.groupme.com/v3/bots/post", "application/json", bytes.NewBuffer(bytesRepresentation))
			}
			c.JSON(http.StatusOK, nil)
		}
	})

	router.Run(":" + port)
}
