package main

import (
	"api/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  string `json:"sex"`
	Id   int    `json:"id"`
}

func GetUserById(res http.ResponseWriter, req *http.Request) {
	url := chi.URLParam(req, "id")
	query := "select id, name, age, sex from users where id = ?"
	rows, err := DB.Query(query, url)
	utils.HttpResponseError(res, err, http.StatusBadRequest)
	defer rows.Close()

	var user User
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Sex)
		utils.HttpResponseError(res, err, http.StatusInternalServerError)
	}
	err = rows.Err()
	utils.HttpResponseError(res, err, http.StatusInternalServerError)

	response := struct {
		User    User   `json:"user"`
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{
		user, true, "User fetched sucessfully!",
	}

	err = json.NewEncoder(res).Encode(response)
	utils.HttpResponseError(res, err, http.StatusInternalServerError)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}

func GetAllUsers(res http.ResponseWriter, req *http.Request) {
	rows, err := DB.Query("select id, name, age, sex from users")
	utils.PanicOnError(err)
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Sex)
		utils.PanicOnError(err)
		users = append(users, user)
	}

	err = rows.Err()
	utils.PanicOnError(err)

	// encode the data as json using encoder since using marshal we send the data in form of []bytes and can't be parsed using JSON.parse() on frontend, rather we need to parse it using TextDecoder()
	response := struct {
		Users   []User `json:"users"`
		Message string `json:"message"`
		Success bool   `json:"success"`
	}{
		users, "Successful", true,
	}
	err = json.NewEncoder(res).Encode(response)
	utils.HttpResponseError(res, err, http.StatusInternalServerError, "JSON encoding error!")

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}

func AddUser(res http.ResponseWriter, req *http.Request) {
	var user User
	err := json.NewDecoder(req.Body).Decode(&user)
	utils.HttpResponseError(res, err, http.StatusBadRequest, "Invalid body params/ JSON parsing error!")

	query, err := DB.Prepare("INSERT INTO users(name, age, sex) VALUES(?, ?, ?)")
	utils.HttpResponseError(res, err, http.StatusInternalServerError, "Error occurred preparing the query!")
	defer query.Close()

	_, err = query.Exec(user.Name, user.Age, user.Sex)
	utils.HttpResponseError(res, err, http.StatusInternalServerError, "Error occurred adding the user to database!")

	rows, err := DB.Query("SELECT id, name, age, sex FROM users")
	utils.HttpResponseError(res, err, http.StatusInternalServerError, "Error occurred fetching users from database!")

	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Sex)

		utils.HttpResponseError(res, err, http.StatusInternalServerError, "Error occurred scanning user row!")

		users = append(users, user)
	}

	response := struct {
		Users   []User `json:"users"`
		Message string `json:"message"`
		Success bool   `json:"success"`
	}{
		users, "Successful", true,
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(response)

	utils.HttpResponseError(res, err, http.StatusInternalServerError, "JSON encoding error!")

}

func UpdateUserById(res http.ResponseWriter, req *http.Request) {

}

func DeleteUserById(res http.ResponseWriter, req *http.Request) {

}

// https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/
