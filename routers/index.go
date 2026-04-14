package routers

import (
	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	// Auth APIs
	authRoutes(api)

	// User APIs
	userRoutes(api)

	// File APIs
	fileRoutes(api)

	// Share APIs
	shareRoutes(api)

	return r

}
