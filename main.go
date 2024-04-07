package main

import (
	"database/sql"
	"fmt"
	"log"
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

	usersRouter := NewUsersRouter()
	router.Get("/", IndexHandler)
	router.Mount("/users", usersRouter)

	server := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%v", PORT),
	}

	fmt.Println("Starting the server at port : ", PORT)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func IndexHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("This is the index page!")
}
