package routers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gorilla/mux"
	"github.com/mayukh551/cloudbox/controllers"
	"github.com/mayukh551/cloudbox/middlewares"
)

func Router() *mux.Router {

	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	// Auth APIs
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", controllers.Login).Methods("POST")
	auth.HandleFunc("/sign-up", controllers.SignUp).Methods("POST")

	// User APIs
	user := api.PathPrefix("/user").Subrouter()
	user.Use(middlewares.Authenticate)
	user.HandleFunc("/get", controllers.GetUserDetails).Methods("GET")
	user.HandleFunc("/update", controllers.UpdateUserDetails).Methods("PUT")
	user.HandleFunc("/delete", controllers.DeleteUser).Methods("DELETE")

	// *************************************

	// File APIs
	file := api.PathPrefix("/file").Subrouter()
	file.Use(middlewares.Authenticate)

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-southeast-2"),
	)

	if err != nil {
		panic(fmt.Errorf("failed to load AWS config, %w", err))
	}

	svc := s3.NewFromConfig(cfg)
	h := controllers.NewHandler(svc)

	file.HandleFunc("/get-list", h.GetList).Methods("GET")
	file.HandleFunc("/download/{type}", h.DownloadFile).Methods("PUT")
	file.HandleFunc("/upload/{type}", h.UploadFile).Methods("POST")
	file.HandleFunc("/delete", h.DeleteFile).Methods("DELETE")

	return r

}
