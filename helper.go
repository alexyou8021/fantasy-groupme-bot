package main

import (
	"log"
	"os"

	"database/sql"
	_ "github.com/lib/pq"
)

func createPlayersTable() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	} else {
                result, err := db.Exec("DROP DATABASE players;")
                result, err = db.Exec("CREATE TABLE players (id int, name varchar(255), position varchar(255));")
		log.Println(result)
		log.Println(err)
	}
}

func storePlayers() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	} else {
                result, err := db.Exec("CREATE TABLE players (id varchar(255), name varchar(255))")
		log.Println(result)
		log.Println(err)
	}
}
