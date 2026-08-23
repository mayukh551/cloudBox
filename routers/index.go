package routers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode("working!")
	})

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
