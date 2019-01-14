package main

import (
	"encoding/json"
	"log"
	"net/http"

	"golang-mvc-boilerplate/server/modules/users"
	shared "golang-mvc-boilerplate/server/sharedVariables"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//App is used to start the app
type App struct {
	Router    *mux.Router
	DBAddress interface{}
	DBName    interface{}
	Port      interface{}
}

func main() {
	var a App
	a.Run()
	a.Serve()
}

//Run is used to create routers
func (a *App) Run() {
	a.Router = users.NewRouter() // create routes
	a.DBAddress = shared.Address
	a.DBName = shared.DbName
	a.Port = shared.Port
}

//Serve is used to serve the routers created
func (a *App) Serve() {
	// launch server with CORS validations
	// these two lines are important in order to allow access from the front-end side to the methods
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	log.Fatal(http.ListenAndServe(":"+string(a.Port.(json.Number)),
		handlers.CORS(allowedOrigins, allowedMethods)(a.Router)))
}
