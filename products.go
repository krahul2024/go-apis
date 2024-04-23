package main

import (
	"api/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

/*
-----------------------------------------------------------------------------------------------------------------
| id          | int           | NO   | PRI | NULL              | auto_increment                                |
| name        | varchar(200)  | YES  |     | NA                |                                               |
| summary     | varchar(500)  | YES  |     | NA                |                                               |
| description | varchar(2000) | YES  |     | NA                |                                               |
| price       | int           | NO   |     | 0                 |                                               |
| quantity    | int           | YES  |     | 0                 |                                               |
| brandId     | int           | NO   | MUL | -1                |                                               |
| categoryId  | int           | NO   | MUL | -1                |                                               |
| imageUrl    | varchar(1000) | YES  |     |                   |                                               |
| createdAt   | datetime      | YES  |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED                             |
| updatedAt   | datetime      | YES  |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED on update CURRENT_TIMESTAMP |
-----------------------------------------------------------------------------------------------------------------

*/

type Product struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	Price       int      `json:"price"`
	Quantity    int      `json:"quantity"`
	Brand       Brand    `json:"brand"`
	Category    Category `json:"category"`
	ImageUrl    string   `json:"imageUrl"`
	ProductCode string   `json:"productCode"`
}

func GetProductById(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	statement := `
		select p.id, p.name, p.summary, p.description, p.price, p.quantity, p.productCode, 
		b.id as brand_id, b.name as brand_name, b.description as brand_description, b.imageUrl as brand_imageUrl, 
		c.id as category_id, c.name as category_name, c.description as category_description, c.imageUrl as category_imageUrl
		from products p
		join brands b on p.brandId = b.id 
		join categories c on p.categoryId = c.id 
		where p.id = ? 
	`
	rows, err := DB.Query(statement, id)
	if err != nil {
		utils.HttpResponseError(res, err, 500, nil)
		return
	}
	defer rows.Close()

	var product Product
	for rows.Next() {
		err = rows.Scan(
			&product.Id,
			&product.Name,
			&product.Summary,
			&product.Description,
			&product.Price,
			&product.Quantity,
			&product.ProductCode,
			&product.Brand.Id,
			&product.Brand.Name,
			&product.Brand.Description,
			&product.Brand.ImageUrl,
			&product.Category.Id,
			&product.Category.Name,
			&product.Category.Description,
			&product.Category.ImageUrl,
		)
		if err != nil {
			utils.HttpResponseError(res, err, 500, nil)
			return
		}
	}
	err = rows.Err()
	if err != nil {
		utils.HttpResponseError(res, err, 500, nil)
		return
	}
	response := struct {
		Product *Product `json:"product"`
		Message string   `json:"message"`
	}{
		&product, "Product Details fetched successfully!",
	}
	if product.Id == 0 {
		response.Product = nil
		response.Message = "No product found!"
	}

	utils.HttpResponseJson(res, response, 200)
}

func AddProduct(res http.ResponseWriter, req *http.Request) {
	var product Product
	err := json.NewDecoder(req.Body).Decode(&product)
	if err != nil {
		utils.HttpResponseError(res, err, 500, nil)
		return
	}
	txn, err := DB.Begin()
	if err != nil {
		utils.HttpResponseError(res, err, 500, nil)
		return
	}
	statement := "insert into products (name, summary, description, price, quantity, brandId, categoryId, imageUrl, productCode) values(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	lastRowId, _, err := utils.DBInsertUpdateDeleteHelper(res, txn, statement, product.Name, product.Summary, product.Description, product.Price, product.Quantity, product.Brand.Id, product.Category.Id, product.ImageUrl, product.ProductCode)
	if err != nil {
		utils.HttpResponseError(res, err, 500, txn)
		return
	}
	query := `
		select p.id, p.name, p.summary, p.description, p.price, p.quantity, p.productCode, 
		b.id as brand_id, b.name as brand_name, b.description as brand_description, b.imageUrl as brand_imageUrl, 
		c.id as category_id, c.name as category_name, c.description as category_description, c.imageUrl as category_imageUrl
		from products p 
		join brands b on p.brandId = b.id 
		join categories c on p.categoryId = c.id 
		where p.id = ? 
		limit 10 
	`
	rows, err := txn.Query(query, lastRowId)
	if err != nil {
		utils.HttpResponseError(res, err, 500, txn)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&product.Id,
			&product.Name,
			&product.Summary,
			&product.Description,
			&product.Price,
			&product.Quantity,
			&product.ProductCode,
			&product.Brand.Id,
			&product.Brand.Name,
			&product.Brand.Description,
			&product.Brand.ImageUrl,
			&product.Category.Id,
			&product.Category.Name,
			&product.Category.Description,
			&product.Category.ImageUrl,
		)
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
		Message string   `json:"message"`
		Product *Product `json:"product"`
	}{
		"Product Added Successfully!", &product,
	}

	utils.HttpResponseJson(res, response, 200)

}

func GetAllProducts(res http.ResponseWriter, req *http.Request) {
	statement := `
		select p.id, p.name, p.summary, p.description, p.price, p.quantity, p.productCode,
		b.id as brand_id, b.name as brand_name, b.description as brand_description, b.imageUrl as brand_imageUrl, 
		c.id as category_id, c.name as category_name, c.description as category_description, c.imageUrl as category_imageUrl
		from products p 
		join brands b on p.brandId = b.id 
		join categories c on p.categoryId = c.id 
		limit 10 
	`
	txn, err := DB.Begin()
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, nil)
		return
	}
	rows, err := txn.Query(statement)
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
		return
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var product Product
		err = rows.Scan(
			&product.Id,
			&product.Name,
			&product.Summary,
			&product.Description,
			&product.Price,
			&product.Quantity,
			&product.ProductCode,
			&product.Brand.Id,
			&product.Brand.Name,
			&product.Brand.Description,
			&product.Brand.ImageUrl,
			&product.Category.Id,
			&product.Category.Name,
			&product.Category.Description,
			&product.Category.ImageUrl,
		)
		if err != nil {
			utils.HttpResponseError(res, err, http.StatusInternalServerError, txn)
			return
		}
		products = append(products, &product)
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
		Total    int        `json:"total"`
		Products []*Product `json:"products"`
		Message  string     `json:"message"`
	}{
		len(products), products, "To get the next batch, use this link!",
	}

	utils.HttpResponseJson(res, response, http.StatusOK)

}

func AddBulkProducts(res http.ResponseWriter, req *http.Request) {

}

func UpdateProduct(res http.ResponseWriter, req *http.Request) {

}

func DeleteProduct(res http.ResponseWriter, req *http.Request) {

}
