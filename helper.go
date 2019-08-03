package main

import (
	"log"
	"os"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"strings"

	"database/sql"
	_ "github.com/lib/pq"
)

func createPlayersTable() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	} else {
		result, err := db.Exec("DROP TABLE players;")
		result, err = db.Exec("CREATE TABLE players (id int, name varchar(255), position varchar(255));")
		log.Println(result)
		log.Println(err)
	}
}

func storePlayers() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	url := "https://api.sleeper.app/v1/players/nfl/"
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	var players map[int]map[string]interface{}
	json.Unmarshal(bodyBytes, &players)

	for _, value := range players {
		id, _ := value["player_id"].(string)
		name, _ := value["full_name"].(string)
		name = strings.Replace(name, "'", "", 1)
		position, _ := value["position"].(string)
		log.Println(name + " " + position + " " + id)
	        _, err := db.Exec("INSERT INTO players VALUES (" + id + ", '" + name + "', '" + position + "');")
                if err != nil {
			log.Fatal(err)
			break
		}
	}
}

func queryPlayer(name string) map[string]interface{} {	
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	result, _ := db.Exec("SELECT * FROM players WHERE name='" + name  + "';")
	log.Println(result)
	return nil
}
