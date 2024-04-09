package main

import (
	"api/utils"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func main() {
	SetEnvs() // Set environment variables first
	InitDB()  // Initialize the database connection
	defer DB.Close()

	PORT := os.Getenv("PORT") // Fetch PORT after setting environment variables

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// User routes
	router.Route("/users", func(r chi.Router) {
		r.Get("/", GetAllUsers)
		r.Post("/", AddUser)
		r.Get("/{id}", GetUserById)
		r.Put("/{id}", UpdateUserById)
		r.Delete("/{id}", DeleteUserById)
	})

	// Auth routes
	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", LoginHandler)
		r.Post("/register", RegistrationHandler)
		r.Get("/logout", LogoutHandler)
	})

	// Brands route
	router.Route("/brands", func(r chi.Router) {
		r.Get("/", GetAllBrands)
		r.Post("/", AddBrand)
		r.Put("/{id}", UpdateBrand)
		r.Delete("/{id}", DeleteBrand)
	})

	// Categories route
	router.Route("/categories", func(r chi.Router) {
		r.Get("/", GetAllCategories)
		r.Post("/", AddCategory)
		r.Put("/{id}", UpdateCategory)
		r.Delete("/{id}", DeleteCategory)
	})

	server := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%v", PORT),
	}

	fmt.Println("Starting the server at port : ", PORT)
	err := server.ListenAndServe()
	utils.FatalError(err, "There was an error starting the server!")

}

func IndexHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("This is the index page!")
}
