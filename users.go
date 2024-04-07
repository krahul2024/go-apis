package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func NewUsersRouter() *chi.Mux {
	userRouter := chi.NewRouter()

	userRouter.Get("/", getAllUsers)
	userRouter.Post("/", addUser)
	userRouter.Get("/{id}", getUserById)
	userRouter.Put("/{id}", updateUserById)
	userRouter.Delete("/{id}", deleteUserById)

	return userRouter
}

func getAllUsers(res http.ResponseWriter, req *http.Request) {
	fmt.Println("This is get all users route!")
}

func addUser(res http.ResponseWriter, req *http.Request) {

}

func getUserById(res http.ResponseWriter, req *http.Request) {

}

func updateUserById(res http.ResponseWriter, req *http.Request) {

}

func deleteUserById(res http.ResponseWriter, req *http.Request) {

}

// https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/
