package groupme

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func msgHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var botResponse msg
		if c.BindJSON(&botResponse) == nil {
			log.Println(botResponse.Text)

			if botResponse.Text == "!help" {
				message := map[string]interface{}{
					"bot_id": os.Getenv("botid"),
					"text":   "I am your chat bot.\nType `!coin` to flip a coin.",
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
					"text":   result,
				}

				bytesRepresentation, err := json.Marshal(message)
				if err != nil {
					log.Fatalln(err)
				}

				http.Post("https://api.groupme.com/v3/bots/post", "application/json", bytes.NewBuffer(bytesRepresentation))
			}
			c.JSON(http.StatusOK, nil)
		}
	}
}
