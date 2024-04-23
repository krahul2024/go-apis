package main

import (
	"api/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

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

func AddBulkProducts(res http.ResponseWriter, req *http.Request) {

}

// this handles dynamic field updates
func UpdateProduct(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	var (
		data   map[string]interface{}
		keys   []string
		values []interface{}
	)
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, nil)
		return
	}
	fmt.Println("Id : ", id)
	for key, value := range data {
		keys = append(keys, key)
		values = append(values, value)
	}
	fmt.Println("Keys : ", keys, "\nValues : ", values)
	statement := "update products set "
	for key, _ := range data {
		if key == "id" {
			continue
		}
		statement += key + "=?, "
	}
	statement = statement[:len(statement)-2]
	statement += " where id = ? "
	values = append(values, id)

	fmt.Println(statement, values)

	txn, err := DB.Begin()
	if err != nil {
		utils.HttpResponseError(res, err, 500, nil)
		return
	}

	lastRowId, numRows, err := utils.DBInsertUpdateDeleteHelper(res, txn, statement, values...)
	if err != nil {
		utils.HttpResponseError(res, err, 500, txn)
		return
	}

	err = txn.Commit()
	if err != nil {
		utils.HttpResponseError(res, err, 500, txn)
		return
	}

	response := struct {
		Message string `json:"message"`
		Status  bool   `json:"status"`
	}{
		"Product Updated successfully!", true,
	}
	fmt.Println(lastRowId, numRows)
	if numRows == 0 {
		response.Message = "No product found or no valid changes for update!"
		response.Status = false
	}
	utils.HttpResponseJson(res, response, 200)

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
	limitParam := req.URL.Query().Get("limit")
	if limitParam == "" {
		limitParam = "5"
	}
	pageParam := req.URL.Query().Get("page")
	if pageParam == "" {
		pageParam = "0"
	}
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		utils.HttpResponseError(res, err, 500, nil)
		return
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		utils.HttpResponseError(res, err, 500, nil)
		return
	}

	offset := page * limit

	fmt.Println(limit)

	statement := `
		select p.id, p.name, p.summary, p.description, p.price, p.quantity, p.productCode,
		b.id as brand_id, b.name as brand_name, b.description as brand_description, b.imageUrl as brand_imageUrl,
		c.id as category_id, c.name as category_name, c.description as category_description, c.imageUrl as category_imageUrl
		from products p
		join brands b on p.brandId = b.id
		join categories c on p.categoryId = c.id
		limit ? offset ? 
	`
	txn, err := DB.Begin()
	if err != nil {
		utils.HttpResponseError(res, err, http.StatusInternalServerError, nil)
		return
	}
	rows, err := txn.Query(statement, limit, offset)
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
		NextLink string     `json:"nextLink"`
		Products []*Product `json:"products"`
	}{
		len(products), fmt.Sprintf("http://localhost:3300/products?limit=%d&page=%d", limit, page+1), products,
	}

	utils.HttpResponseJson(res, response, http.StatusOK)

}

func DeleteProduct(res http.ResponseWriter, req *http.Request) {

}
