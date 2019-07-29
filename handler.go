package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type msg struct {
	Text string `json:"text"`
}

func sendPost(text string) {
	message := map[string]interface{}{
		"bot_id": os.Getenv("botid"),
		"text":   text,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	http.Post("https://api.groupme.com/v3/bots/post", "application/json", bytes.NewBuffer(bytesRepresentation))
}

func msgHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var botResponse msg
		if c.BindJSON(&botResponse) == nil {
			log.Println(botResponse.Text)

			if botResponse.Text == "!help" {
				sendPost("I am your chat bot.\nType `!coin` to flip a coin.")
			}

			if botResponse.Text == "!coin" {
				result := "Your coin landed on HEADS."
				if rand.Intn(2) == 1 {
					result = "Your coin landed on TAILS."
				}
				sendPost(result)
			}
			c.JSON(http.StatusOK, nil)
		}
	}
}
