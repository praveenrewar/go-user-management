package main

import (
	"log"
	"net/http"
	"strconv"

	"go-user-management/server/modules/users"
	shared "go-user-management/server/sharedVariables"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//App is used to start the app
type App struct {
	Router    *mux.Router
	DBAddress string
	DBName    string
	Port      int
}

//Main function to run and serve the App
func main() {
	var a App
	a.Run()
	a.Serve()
}

//Run is used to create routers
func (a *App) Run() {
	a.Router = mux.NewRouter()
	a.Router = users.AddRouters(a.Router)
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
	allowedHeader := handlers.AllowedHeaders([]string{"Accept", "Authorization", "X-CSRF-Token", "X-Requested-With", "Content-Type"})
	log.Println("Serving application at PORT:" + strconv.Itoa(a.Port))
	log.Fatal(http.ListenAndServe(":"+string(strconv.Itoa(a.Port)),
		handlers.CORS(allowedOrigins, allowedMethods, allowedHeader)(a.Router)))
}
