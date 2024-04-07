package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func InitDB() {
	DB_USER, DB_PORT, DB_PASSWORD := os.Getenv("DB_USER"), os.Getenv("DB_PORT"), os.Getenv("DB_PASSWORD")
	var err error
	DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:%s)/first", DB_USER, DB_PASSWORD, DB_PORT))
	if err != nil {
		log.Fatal("There was an error connecting to the database!\n", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("There was an error connecting to the database!\n", err)
	}

	log.Println("Connected to the database successfully!")
}
