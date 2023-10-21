package routes

import (
	"github.com/gorilla/mux"
	"github.com/shuaibu222/go-bookstore/api/controllers"
	"github.com/shuaibu222/go-bookstore/auth"
)

var UsersRoutes = func(router *mux.Router) {
	router.HandleFunc("/v1/api/users", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/v1/api/users", controllers.GetAllUsers).Methods("GET")
	router.HandleFunc("/v1/api/login", auth.Login).Methods("POST")
	router.HandleFunc("/v1/api/logout", auth.Logout).Methods("POST")
	router.HandleFunc("/v1/api/refresh", auth.Refresh).Methods("GET")
	router.HandleFunc("/v1/api/users/{id}", auth.AuthMiddleware(controllers.GetUserByID)).Methods("GET")
	router.HandleFunc("/v1/api/users/{id}", auth.AuthMiddleware(controllers.UpdateUser)).Methods("PUT")
	router.HandleFunc("/v1/api/users/{id}", auth.AuthMiddleware(controllers.DeleteUserById)).Methods("DELETE")
}
