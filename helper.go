package main

import (
	"log"
	"os"

	"database/sql"
	_ "github.com/lib/pq"
)

func storePlayers() {
	_, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("db")
	}
}
