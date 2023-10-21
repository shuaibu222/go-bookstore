package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/shuaibu222/go-bookstore/api/routes"
	"github.com/shuaibu222/go-bookstore/config"
)

func main() {
	r := mux.NewRouter()

	routes.BooksRoutes(r)
	routes.UsersRoutes(r)
	http.Handle("/", r)

	config, err := config.LoadConfig()
	if err != nil {
		log.Println("Error while loading envs: ", err)
	}

	log.Println("Server started.........")

	log.Fatal(http.ListenAndServe(":"+config.WebPort,
		handlers.CORS(
			handlers.AllowCredentials(),
			handlers.AllowedOrigins([]string{"http://*, https://*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}),
		)(r)),
	)
}
