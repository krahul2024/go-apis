package main

import (
	"api/utils"
	"database/sql"
	"fmt"
	"log"
	"os"
)

func InitDB() {
	DB_USER, DB_PORT, DB_PASSWORD := os.Getenv("DB_USER"), os.Getenv("DB_PORT"), os.Getenv("DB_PASSWORD")
	var err error
	DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:%s)/first", DB_USER, DB_PASSWORD, DB_PORT))
	utils.FatalError(err)

	err = DB.Ping()
	utils.FatalError(err)

	log.Println("Connected to the database successfully!")
}
