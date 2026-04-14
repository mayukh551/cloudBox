package routers

import (
	"github.com/gorilla/mux"
	"github.com/mayukh551/cloudbox/controllers"
	"github.com/mayukh551/cloudbox/middlewares"
)

func shareRoutes(api *mux.Router) {
	share := api.PathPrefix("/share").Subrouter()
	share.Use(middlewares.Authenticate)

	share.HandleFunc("/share", controllers.Share).Methods("POST")
	share.HandleFunc("/get-shared-list", controllers.ListShares).Methods("GET")
	share.HandleFunc("/get-shared-with-me", controllers.ListSharedWithMe).Methods("GET")
}
