package routers

import (
	"github.com/gorilla/mux"
	"github.com/mayukh551/cloudbox/controllers"
	"github.com/mayukh551/cloudbox/middlewares"
	"github.com/mayukh551/cloudbox/utils"
)

func fileRoutes(api *mux.Router) error {
	// File APIs
	file := api.PathPrefix("/file").Subrouter()
	file.Use(middlewares.Authenticate)

	s3, err := utils.LoadAWSConfig()
	if err != nil {
		return err
	}

	h := controllers.NewHandler(s3)

	file.HandleFunc("/get-list", h.GetList).Methods("GET")
	file.HandleFunc("/download/{type}", h.DownloadFile).Methods("PUT")
	file.HandleFunc("/upload/{type}", h.UploadFile).Methods("POST")
	file.HandleFunc("/rename", h.Rename).Methods("PUT")
	file.HandleFunc("/delete", h.DeleteFile).Methods("PUT")

	return nil
}
