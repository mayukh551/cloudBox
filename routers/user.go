package routers

import (
	"github.com/gorilla/mux"
	"github.com/mayukh551/cloudbox/controllers"
	"github.com/mayukh551/cloudbox/middlewares"
)

func userRoutes(api *mux.Router) {

	user := api.PathPrefix("/user").Subrouter()

	user.Use(middlewares.Authenticate)
	user.HandleFunc("/get", controllers.GetUserDetails).Methods("GET")
	user.HandleFunc("/update", controllers.UpdateUserDetails).Methods("PUT")
	user.HandleFunc("/delete", controllers.DeleteUser).Methods("DELETE")
	user.HandleFunc("/find", controllers.FindUserByEmail).Methods("PUT")
}
