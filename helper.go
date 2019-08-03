package main

import (
	"log"
	"os"
	"io/ioutil"
	"encoding/json"
	"net/http"

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
	//db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	//if err != nil {
	//	log.Fatal(err)
	//} else {
	//        result, err := db.Exec("CREATE TABLE players (id varchar(255), name varchar(255))")
	//	log.Println(result)
	//	log.Println(err)
	//}
	url := "https://api.sleeper.app/v1/players/nfl/"
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	var players []map[string]interface{}
	json.Unmarshal(bodyBytes, &players)
	log.Println(players)

}
