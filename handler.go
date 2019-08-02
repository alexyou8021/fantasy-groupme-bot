package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
        "strings"
        "io/ioutil"

	"github.com/gin-gonic/gin"
)

type msg struct {
	Text string `json:"text"`
}

type League struct {
    Response map[string][]map[string]string `json:"response"`
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
                        fields := strings.Fields(botResponse.Text)
                        log.Println(fields)
                        if len(fields) == 0 {
			    c.JSON(http.StatusOK, nil)
                            return
                        }

			if fields[0] == "!help" {
				sendPost("I am your chat bot.\nType `!coin` to flip a coin.\nType `!smack` to trash talk.")
			}

			if fields[0] == "!coin" {
				result := "Your coin landed on HEADS."
				if rand.Intn(2) == 1 {
					result = "Your coin landed on TAILS."
				}
				sendPost(result)
			}

			if fields[0] == "!smack" {
                            groupid := os.Getenv("groupid")
                            url1 := "https://api.groupme.com/v3/groups/" + groupid + "?token="
                            url1 = url1 + os.Getenv("token")
                            resp1, _ := http.Get(url1)

                            defer resp1.Body.Close()
                            bodyBytes1, _ := ioutil.ReadAll(resp1.Body)
                            var league League
                            json.Unmarshal(bodyBytes1, &league)

                            members := league.Response["members"]
                            memberNum := -1

                            for i := 0; i < len(members); i++ {
                                if len(fields) == 1 {
                                    break
                                } else if len(fields) == 2 {
                                    if fields[1] == "@" + members[i]["nickname"] {
                                        memberNum = i
                                        break
                                    }
                                } else {
                                    if fields[1] + " " + fields[2] == "@" + members[i]["nickname"] {
                                        memberNum = i
                                        break
                                    }
                                }
                            }

                            if memberNum == -1 {
                                memberNum = rand.Intn(len(members))
                            }

                            nickname := strings.Replace(members[memberNum]["nickname"], " ", "%20", 1)

                            url2 := "https://insult.mattbas.org/api/insult?who=" + nickname
                            log.Println(url2)
                            resp2, _ := http.Get(url2)

                            defer resp2.Body.Close()
                            bodyBytes2, _ := ioutil.ReadAll(resp2.Body)

                            result := "@" + string(bodyBytes2)
                            sendPost(result)
                        }

                        if fields[0] == "!stats" {
                            if len(fields) <= 3 || len(fields) >= 6 {
			        c.JSON(http.StatusOK, nil)
                                return
                            }
                            name := fields[1] + " " + fields[2]
                            season := fields[3]
                            week := ""
                            if len(fields) == 5 {
                                week = fields[4]
                            }

                            url := "https://api.sleeper.app/v1/stats/nfl/regular/" + season + "/" + week
                            resp, _ := http.Get(url)
                            defer resp.Body.Close()
                            bodyBytes, _ := ioutil.ReadAll(resp.Body)
                            var stats map[int]map[string]float32
                            json.Unmarshal(bodyBytes, &stats)
                            log.Println(url)
                            log.Println(stats[1034])
                            log.Println(name)
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
