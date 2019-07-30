package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

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

func reminderHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("reminders") != "on" {
			c.JSON(http.StatusOK, nil)
			return
		}	
		day := int(time.Now().Weekday())

		if day == 0 {
			sendPost("Reminder:\nSunday games start soon.\nDon't forget to set your lineups!")
		}
		if day == 2 {
			sendPost("Reminder:\nWaivers will be process soon.\nDon't forget to set your waivers!")
		}
		if day == 4 {
			sendPost("Reminder:\nThursday games start soon.\nDon't forget to set your lineups!")
		}
		c.JSON(http.StatusOK, nil)
	}
}
