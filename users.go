package main

import (
	"log"
	"net/http"
)

func GetAllUsers(res http.ResponseWriter, req *http.Request) {
	rows, err := DB.Query("select name, age, gender, country from customers")
	if err != nil {
		log.Printf("There was an error!, \n%v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name, &age, &gender)
		if err != nil {

		}
	}
}

func AddUser(res http.ResponseWriter, req *http.Request) {

}

func GetUserById(res http.ResponseWriter, req *http.Request) {

}

func UpdateUserById(res http.ResponseWriter, req *http.Request) {

}

func DeleteUserById(res http.ResponseWriter, req *http.Request) {

}

// https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/
