package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/handlers"

	"./server/modules/users"
	"./server/sharedVariables"
)

func main() {

	var port = shared.Port
	router := users.NewRouter() // create routes
	// these two lines are important in order to allow access from the front-end side to the methods
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})

	// launch server with CORS validations
	log.Fatal(http.ListenAndServe(":"+string(port.(json.Number)),
		handlers.CORS(allowedOrigins, allowedMethods)(router)))
}
