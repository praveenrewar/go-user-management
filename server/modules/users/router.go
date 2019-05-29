package users

import (
	"net/http"

	"go-user-management/server/logger"
	jwtauthenticate "go-user-management/server/middlewares/jwtAuthenticate"
	usersmiddleware "go-user-management/server/middlewares/usersMiddleware"

	"github.com/gorilla/mux"
)

var controller = &Controller{Repository: Repository{}}

// Route defines a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes defines the list of routes of our API
type Routes []Route

var routes = Routes{
	Route{
		Name:        "GetUsers",
		Method:      "GET",
		Pattern:     "/get-users",
		HandlerFunc: controller.GetUsers,
	},
	Route{
		Name:        "Login",
		Method:      "POST",
		Pattern:     "/login",
		HandlerFunc: controller.Login,
	},
	Route{
		Name:        "Signup",
		Method:      "POST",
		Pattern:     "/signup",
		HandlerFunc: controller.Signup,
	},
	Route{
		Name:        "UpdatePassword",
		Method:      "POST",
		Pattern:     "/update-password",
		HandlerFunc: controller.UpdatePassword,
	},
	Route{
		Name:        "DeleteUser",
		Method:      "DELETE",
		Pattern:     "/delete-user/{user_id}",
		HandlerFunc: controller.DeleteUser,
	}}

//AddRouters configures a new router to the API
func AddRouters(router *mux.Router) *mux.Router {
	// router := mux.NewRouter().StrictSlash(true)

	//signup
	router.Methods()

	for _, route := range routes {
		var handler http.Handler
		// var validator http.Handler
		handler = route.HandlerFunc
		if route.Name == "Signup" {
			handler = usersmiddleware.SignupValidator(handler)
		}
		if route.Name == "Login" {
			handler = usersmiddleware.LoginValidator(handler)
		}
		if route.Name == "UpdatePassword" {
			handler = usersmiddleware.UpdatePasswordValidator(handler)
			handler = jwtauthenticate.Authenticate(handler)
		}
		if route.Name == "GetUsers" {
			handler = jwtauthenticate.Authenticate(handler)
		}
		if route.Name == "DeleteUser" {
			handler = jwtauthenticate.Authenticate(handler)
		}
		handler = logger.Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	return router
}
