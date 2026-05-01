package routers

import (
	"github.com/gorilla/mux"
	"github.com/mayukh551/cloudbox/controllers"
	"github.com/mayukh551/cloudbox/middlewares"
)

func shareRoutes(api *mux.Router) {
	share := api.PathPrefix("/share").Subrouter()
	share.Use(middlewares.Authenticate)

	share.HandleFunc("/create", controllers.Share).Methods("POST")
	share.HandleFunc("/list", controllers.ListShares).Methods("GET")
	share.HandleFunc("/list-with-me", controllers.ListSharedWithMe).Methods("GET")
}
