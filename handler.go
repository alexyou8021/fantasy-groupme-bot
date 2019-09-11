package main

import (
	"bytes"
	"sort"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type msg struct {
	Text       string `json:"text"`
	GroupId    string `json:"group_id"`
	Id         string `json:"id"`
	Name       string `json:"name"`
	SenderId   string `json:"sender_id"`
	SenderType string `json:"sender_type"`
	UserId     string `json:"user_id"`
}

type League struct {
	Response map[string][]map[string]string `json:"response"`
}

type Bot struct {
	Name      string `json:"name"`
	BotId     string `json:"bot_id"`
	GroupId   string `json:"group_id"`
	GroupName string `json:"group_name"`
}

type Team struct {
	Name string
	Wins float64
	Losses float64
	Waiver float64
	Budget float64
}

func sendPost(text string, bot_id string) {
	message := map[string]interface{}{
		"bot_id": bot_id,
		"text":   text,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	http.Post("https://api.groupme.com/v3/bots/post", "application/json", bytes.NewBuffer(bytesRepresentation))
}

func getBots() []Bot {
	url := "https://api.groupme.com/v3/bots?token=" + os.Getenv("token")
	log.Println(url)
	resp, _ := http.Get(url)
	log.Println(resp)
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	log.Println(bodyBytes)
	var response map[string]interface{}
	json.Unmarshal(bodyBytes, &response)
	dict, _ := json.Marshal(response["response"])
	var bots []Bot
	json.Unmarshal(dict, &bots)
	log.Println(bots)
	return bots
}

func msgHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var botResponse msg
		if c.BindJSON(&botResponse) == nil {
			bots := getBots()
			fields := strings.Fields(botResponse.Text)
			groupId := botResponse.GroupId
			botId := ""

			for _, bot := range bots {
				if bot.GroupId == groupId {
					botId = bot.BotId
					break
				}
			}

			log.Println(botId)
			log.Println(fields)

			if len(fields) == 0 {
				c.JSON(http.StatusOK, nil)
				return
			}

			if fields[0] == "!help" {
				if botResponse.GroupId == os.Getenv("htown") {
					sendPost("I am your chat bot.\nType `!coin` to flip a coin.\nType `!smack @someone` to talk trash.\nType `!stats player season week` for stats.\nType `!draft` for draft info.\nType `!standings` for league standings.", botId)
				} else {
					sendPost("I am your chat bot.\nType `!coin` to flip a coin.\nType `!smack @someone` to talk trash.\nType `!stats player season week` for stats.", botId)
				}
			} else if fields[0] == "!coin" {
				result := "Your coin landed on HEADS."
				if rand.Intn(2) == 1 {
					result = "Your coin landed on TAILS."
				}
				sendPost(result, botId)
			} else if fields[0] == "!draft" {
				message := os.Getenv("draft")
				sendPost(message, botId)
			} else if fields[0] == "!smack" {
				groupid := botResponse.GroupId
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
						if fields[1] == "@"+members[i]["nickname"] {
							memberNum = i
							break
						}
					} else {
						if fields[1]+" "+fields[2] == "@"+members[i]["nickname"] {
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
				sendPost(result, botId)
			} else if fields[0] == "!standings" {
				league := os.Getenv("league")

				url1 := "https://api.sleeper.app/v1/league/" + league + "/users"
				resp1, _ := http.Get(url1)

				defer resp1.Body.Close()
				bodyBytes1, _ := ioutil.ReadAll(resp1.Body)
				var users []map[string]interface{}
				json.Unmarshal(bodyBytes1, &users)

				usernames := make(map[string]string)
				for _, value := range users {
					id, _ := value["user_id"].(string)
					display_name, _ := value["display_name"].(string)
					usernames[id] = display_name
				}

				log.Println(usernames)

				url2 := "https://api.sleeper.app/v1/league/" + league + "/rosters"
				resp2, _ := http.Get(url2)

				defer resp2.Body.Close()
				bodyBytes2, _ := ioutil.ReadAll(resp2.Body)
				var rosters []map[string]interface{}
				json.Unmarshal(bodyBytes2, &rosters)
				standings := make([]Team, 12)
				for key, value := range rosters {
					owner_id, _ := value["owner_id"].(string)
					display_name := usernames[owner_id]
					settings := value["settings"].(map[string]interface{})
					var team Team
					team.Name = display_name
					team.Wins, _ = settings["wins"].(float64)
					team.Losses, _ = settings["losses"].(float64)
					team.Waiver, _ = settings["waiver_position"].(float64)
					team.Budget, _ = 200 - settings["waiver_budget_used"].(float64)
					standings[key] = team
				}
				
				var teamList []Team
				for _, value := range standings {
					teamList = append(teamList, value)
				}

				sort.Slice(teamList, func(i, j int) bool {
					if teamList[i].Wins == teamList[j].Wins {
        					return teamList[i].Waiver < teamList[j].Waiver
					}
        				return teamList[i].Wins > teamList[j].Wins
    				})

				message := "Name      Record Waiver\n-----------------------------\n"
				for _, value := range teamList {
					message = message + value.Name + "\n"
					message = fmt.Sprintf("%s                   %0.f-%0.f      %0.f\n", message, value.Wins, value.Losses, value.Budget)
				}

				sendPost(message, botId)
			} else if fields[0] == "!stats" {
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

				player, err := queryPlayer(name)
				if err != nil {
					return
				}
				if player.Name == "" {
					sendPost("Player Not Found.", botId)
					return
				}

				url := "https://api.sleeper.app/v1/stats/nfl/regular/" + season + "/" + week
				resp, _ := http.Get(url)
				defer resp.Body.Close()
				bodyBytes, _ := ioutil.ReadAll(resp.Body)
				var stats map[int]map[string]float32
				json.Unmarshal(bodyBytes, &stats)
				stat := stats[player.Id]

				log.Println(url)
				log.Println(stat)
				log.Println(player.Name)

				pts := fmt.Sprintf("%.1f", stat["pts_half_ppr"])
				message := player.Name + ": " + pts + " pts\n"
				if player.Position == "WR" || player.Position == "TE" {
					rec_tgt := fmt.Sprintf("%.0f", stat["rec_tgt"])
					rec := fmt.Sprintf("%.0f", stat["rec"])
					rec_yd := fmt.Sprintf("%.0f", stat["rec_yd"])
					rec_td := fmt.Sprintf("%.0f", stat["rec_td"])
					message = message + "- Targets: " + rec_tgt + "\n"
					message = message + "- Catches: " + rec + "\n"
					message = message + "- Yards: " + rec_yd + "\n"
					message = message + "- TDs: " + rec_td + "\n"
					sendPost(message, botId)
				} else if player.Position == "RB" {
					rush_att := fmt.Sprintf("%.0f", stat["rush_att"])
					rush_yd := fmt.Sprintf("%.0f", stat["rush_yd"])
					rush_td := fmt.Sprintf("%.0f", stat["rush_td"])
					rec_tgt := fmt.Sprintf("%.0f", stat["rec_tgt"])
					rec := fmt.Sprintf("%.0f", stat["rec"])
					rec_yd := fmt.Sprintf("%.0f", stat["rec_yd"])
					rec_td := fmt.Sprintf("%.0f", stat["rec_td"])
					message = message + "- Rush Att: " + rush_att + "\n"
					message = message + "- Rush Yards: " + rush_yd + "\n"
					message = message + "- Rush TDs: " + rush_td + "\n"
					message = message + "- Targets: " + rec_tgt + "\n"
					message = message + "- Catches: " + rec + "\n"
					message = message + "- Rec Yards: " + rec_yd + "\n"
					message = message + "- Rec TDs: " + rec_td + "\n"
					sendPost(message, botId)
				} else if player.Position == "QB" {
					sendPost(player.Name+": "+pts+" pts", botId)
					pass_yd := fmt.Sprintf("%.0f", stat["pass_yd"])
					pass_td := fmt.Sprintf("%.0f", stat["pass_td"])
					pass_int := fmt.Sprintf("%.0f", stat["pass_int"])
					rush_yd := fmt.Sprintf("%.0f", stat["rush_yd"])
					rush_td := fmt.Sprintf("%.0f", stat["rush_td"])
					fum_lost := fmt.Sprintf("%.0f", stat["fum_lost"])
					message = message + "- Passing Yards: " + pass_yd + "\n"
					message = message + "- Passing TDs: " + pass_td + "\n"
					message = message + "- Passing INTs: " + pass_int + "\n"
					message = message + "- Rush Yards: " + rush_yd + "\n"
					message = message + "- Rush TDs: " + rush_td + "\n"
					message = message + "- Fumbles: " + fum_lost + "\n"
				} else {
					sendPost(player.Name+": "+pts+" pts", botId)
				}
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
			sendPost("Reminder:\nSunday games start soon.\nDon't forget to set your lineups!", os.Getenv("botid"))
		}
		if day == 2 {
			sendPost("Reminder:\nWaivers will be process soon.\nDon't forget to set your waivers!", os.Getenv("botid"))
		}
		if day == 4 {
			sendPost("Reminder:\nThursday games start soon.\nDon't forget to set your lineups!", os.Getenv("botid"))
		}
		c.JSON(http.StatusOK, nil)
	}
}
