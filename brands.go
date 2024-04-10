package main

import (
	"api/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Brand struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageUrl    string `json:"imageUrl"`
}

func GetAllBrands(res http.ResponseWriter, req *http.Request) {
	txn, err := DB.Begin()
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, nil)
		return
	}

	statement := "Select id, name, description, imageUrl from brands"
	rows, err := txn.Query(statement)
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}
	defer rows.Close()

	var brands []Brand

	for rows.Next() {
		var brand Brand
		err := rows.Scan(&brand.Id, &brand.Name, &brand.Description, &brand.ImageUrl)
		if err != nil {
			utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
			return
		}
		brands = append(brands, brand)
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
		Brands  []Brand `json:"brands"`
		Status  bool    `json:"status"`
		Message string  `json:"message"`
	}{
		brands, true, "Brands fetched successfully!",
	}

	utils.HttpResponseJson(res, response, 200)
}

func AddBrand(res http.ResponseWriter, req *http.Request) {
	var brand Brand
	err := json.NewDecoder(req.Body).Decode(&brand)
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusBadRequest, nil)
		return
	}
	txn, err := DB.Begin()
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusBadRequest, nil)
		return
	}

	statement := "insert into brands (name, description, imageUrl) values(?, ?, ?)"
	_, _, err = utils.DBInsertUpdateDeleteHelper(res, txn, statement, brand.Name, brand.Description, brand.ImageUrl)
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}

	// get all the brands
	statement = "select id, name, description, imageUrl from brands"
	rows, err := txn.Query(statement)
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}
	defer rows.Close()

	var brands []Brand
	for rows.Next() {
		var tempBrand Brand
		err := rows.Scan(&tempBrand.Id, &tempBrand.Name, &tempBrand.Description, &tempBrand.ImageUrl)
		if err != nil {
			utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
			return
		}
		brands = append(brands, tempBrand)
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
		Brands  []Brand `json:"brands"`
		Message string  `json:"message"`
	}{
		brands, "Brand Added successfully!",
	}

	utils.HttpResponseJson(res, response, http.StatusAccepted)
}

func UpdateBrand(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	var brand Brand
	err := json.NewDecoder(req.Body).Decode(&brand)
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusBadRequest, nil)
		return
	}
	statement := "update brands set name=?, description=?, imageUrl=? where id=?"
	txn, err := DB.Begin()
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, nil)
		return
	}

	_, numRows, err := utils.DBInsertUpdateDeleteHelper(res, txn, statement, brand.Name, brand.Description, brand.ImageUrl, id)
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}

	statement = "select id, name, description, imageUrl from brands where id=?"
	rows, err := txn.Query(statement, id)
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}

	var updatedBrand Brand
	for rows.Next() {
		err := rows.Scan(&updatedBrand.Id, &updatedBrand.Name, &updatedBrand.Description, &updatedBrand.ImageUrl)
		if err != nil {
			utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
			return
		}
	}

	err = txn.Commit()
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}

	response := struct {
		Brand   Brand  `json:"brand"`
		Message string `json:"message"`
	}{
		updatedBrand, "Brand updated successfull!",
	}

	if numRows == 0 {
		response.Message = "No brand found!"
		utils.HttpResponseJson(res, response, http.StatusNotFound)
		return
	}

	utils.HttpResponseJson(res, response, http.StatusAccepted)
}

func DeleteBrand(res http.ResponseWriter, req *http.Request) {

}
