package main

import (
	"api/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Sex      string `json:"sex"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	City     string `json:"city"`
	Country  string `json:"country"`
}

func DeleteUserById(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	txn, err := DB.Begin()
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, nil)
		return
	}
	statement := "delete from users where id=?"
	_, numRows, err := utils.DBInsertUpdateDeleteHelper(res, txn, statement, id)
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}

	err = txn.Commit()
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}

	response := struct {
		Message string `json:"message"`
		Status  bool   `json:"status"`
	}{"User deleted successfully!", true}

	if numRows == 0 {
		response.Message = "No user found with this id!"
		response.Status = false
	}

	utils.HttpResponseJson(res, response, http.StatusOK)
}

func UpdateUserById(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	var user User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusBadRequest, nil)
		return
	}

	txn, err := DB.Begin()
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusBadRequest, nil)
		return
	}

	statement := "update users set name=?, age=?, sex=? where id=?"
	_, numRows, err := utils.DBInsertUpdateDeleteHelper(res, txn, statement, id)
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}

	// since the user details are updated , we have to return the updated user information
	rows, err := txn.Query("select id, name, age, sex from users where id = ?", id)
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Sex)
		if err != nil {
			utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
			return
		}
	}
	err = rows.Err()
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}

	err = txn.Commit()
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}

	response := struct {
		User    User   `json:"user"`
		Message string `json:"message"`
		Status  bool   `json:"status"`
	}{
		user, "User information updated successfully!", true,
	}

	if numRows == 0 {
		response.Message = "No user found with this id!"
		response.Status = false
		response.User = User{}
	}

	utils.HttpResponseJson(res, response, 200)
}

func GetUserById(res http.ResponseWriter, req *http.Request) {
	url := chi.URLParam(req, "id")

	txn, err := DB.Begin()
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusBadRequest, nil)
		return
	}
	query := "select id, name, age, sex, username, email, phone, password, city, country from users where id = ?"
	rows, err := txn.Query(query, url)
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}
	defer rows.Close()

	var user User
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Sex, &user.Username, &user.Email, &user.Phone, &user.Password, &user.City, &user.Country)
		if err != nil {
			utils.HttpResponseError(res, err, http.StatusBadRequest, txn)
			return
		}
	}
	err = rows.Err()
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}

	err = txn.Commit()
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}

	response := struct {
		User    User   `json:"user"`
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}{
		user, true, "User fetched sucessfully!",
	}

	if user.Id == 0 {
		response.Message = "No user found!"
		response.Status = true
		response.User = User{}
	}
	utils.HttpResponseJson(res, response, 200)
}

func GetAllUsers(res http.ResponseWriter, req *http.Request) {
	txn, err := DB.Begin()
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusBadRequest, nil)
		return
	}
	rows, err := txn.Query("SELECT id, name, age, sex, username, email, phone, password, city, country FROM users")
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusBadRequest, txn)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Sex, &user.Username, &user.Email, &user.Phone, &user.Password, &user.City, &user.Country)
		if err != nil {
			utils.HttpResponseError(res, err, http.StatusBadRequest, txn)
			return
		}
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusBadRequest, txn)
		return
	}

	err = txn.Commit()
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}

	// encode the data as json using encoder since using marshal we send the data in form of []bytes and can't be parsed using JSON.parse() on frontend, rather we need to parse it using TextDecoder()
	response := struct {
		Users   []User `json:"users"`
		Message string `json:"message"`
		Status  bool   `json:"status"`
	}{
		users, "Successful", true,
	}
	utils.HttpResponseJson(res, response, 200)
}

func AddUser(res http.ResponseWriter, req *http.Request) {
	var newUser User
	err := json.NewDecoder(req.Body).Decode(&newUser)
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, nil)
		return
	}

	// Transaction starts from here
	txn, err := DB.Begin()
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, nil)
		return
	}

	statement := "INSERT INTO users(name, age, sex, username, email, phone, password, city, country) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	lastRow, numRows, err := utils.DBInsertUpdateDeleteHelper(res, txn, statement, newUser.Name, newUser.Age, newUser.Sex, newUser.Username, newUser.Email, newUser.Phone, newUser.Password, newUser.City, newUser.Country)
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}

	rows, err := txn.Query("SELECT id, name, age, sex, username, email, phone, password, city, country FROM users")
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Sex, &user.Username, &user.Email, &user.Phone, &user.Password, &user.City, &user.Country)
		if err != nil {
			utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
			return
		}
		users = append(users, user)
	}

	response := struct {
		Users   []User `json:"users"`
		Message string `json:"message"`
		Status  bool   `json:"status"`
	}{
		users, "Successful", true,
	}

	err = txn.Commit()
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}

	if numRows != 0 || lastRow != 0 {
		response.Message = "User added successfully!"
	}

	utils.HttpResponseJson(res, response, 200)
}

// https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/
