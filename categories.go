package main

import (
	"api/utils"
	"encoding/json"
	"net/http"
)

type Category struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageUrl    string `json:"imageUrl"`
}

func GetAllCategories(res http.ResponseWriter, req *http.Request) {
	statement := "select id, name, description, imageUrl from categories"
	rows, err := DB.Query(statement)
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, nil)
		return
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.Id, &category.Name, &category.Description, &category.ImageUrl)
		if err != nil {
			utils.HttpResponseError(res, err, http.StatusInternalServerError, nil)
			return
		}
		categories = append(categories, category)
	}
	if err = rows.Err(); err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, nil)
		return
	}

	response := struct {
		Total      int        `json:"total"`
		Message    string     `json:"message"`
		Categories []Category `json:"categories"`
	}{
		len(categories), "Categories fetched successfully!", categories,
	}
	utils.HttpResponseJson(res, response, http.StatusOK)
}

func AddCategory(res http.ResponseWriter, req *http.Request) {
	var category Category
	err := json.NewDecoder(req.Body).Decode(&category)
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, nil)
		return
	}

	txn, err := DB.Begin()
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, nil)
		return
	}

	statement := "insert into categories (name, description, imageUrl) values (?, ?, ?)"
	lastRow, _, err := utils.DBInsertUpdateDeleteHelper(res, txn, statement, category.Name, category.Description, category.ImageUrl)
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}

	statement = "select id, name, description, imageUrl from categories where id=?"
	rows, err := txn.Query(statement, lastRow)
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}
	defer rows.Close()

	var addedCategory Category
	for rows.Next() {
		err := rows.Scan(&addedCategory.Id, &addedCategory.Name, &addedCategory.Description, &addedCategory.ImageUrl)
		if err != nil {
			utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
			return
		}
	}

	if err = rows.Err(); err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}

	if err = txn.Commit(); err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}

	response := struct {
		Category Category `json:"category"`
		Message  string   `json:"message"`
	}{
		addedCategory, "Category added successfully!",
	}

	utils.HttpResponseJson(res, response, http.StatusAccepted)
}

func UpdateCategory(res http.ResponseWriter, req *http.Request) {

}

func DeleteCategory(res http.ResponseWriter, req *http.Request) {

}
