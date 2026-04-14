package routers

import (
	"github.com/gorilla/mux"
	"github.com/mayukh551/cloudbox/controllers"
)

func authRoutes(api *mux.Router) {
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", controllers.Login).Methods("POST")
	auth.HandleFunc("/sign-up", controllers.SignUp).Methods("POST")
}
