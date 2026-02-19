package routers

import (
	"github.com/gorilla/mux"
	"github.com/mayukh551/cloudbox/controllers"
	"github.com/mayukh551/cloudbox/middlewares"
)

func Router() *mux.Router {

	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	// auth routes
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", controllers.Login).Methods("POST")
	auth.HandleFunc("/sign-up", controllers.SignUp).Methods("POST")

	// user routes
	user := api.PathPrefix("/user").Subrouter()
	user.Use(middlewares.Authenticate)
	user.HandleFunc("/get", controllers.GetUserDetails).Methods("GET")
	user.HandleFunc("/update", controllers.UpdateUserDetails).Methods("PUT")
	user.HandleFunc("/delete", controllers.DeleteUser).Methods("DELETE")

	file := api.PathPrefix("/file").Subrouter()
	file.Use(middlewares.Authenticate)

	file.HandleFunc("/download/{type}", controllers.DownloadFile).Methods("PUT") // single or multiple uploads
	file.HandleFunc("/upload/{type}", controllers.UploadFile).Methods("POST")    // single or multiple uploads
	// file.HandleFunc("/share/{type}")    // type: public, invited only via email
	// file.HandleFunc("/search")          // search files

	return r

}
