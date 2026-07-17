package routers

import (
	"github.com/gorilla/mux"
	c "github.com/mayukh551/cloudbox/controllers"
	m "github.com/mayukh551/cloudbox/middlewares"
	"github.com/mayukh551/cloudbox/utils"
)

func fileRoutes(api *mux.Router) error {
	// File APIs
	file := api.PathPrefix("/file").Subrouter()
	file.Use(m.Authenticate)

	s3, err := utils.LoadAWSConfig()
	if err != nil {
		return err
	}

	h := c.NewHandler(s3)

	// file.HandleFunc("/get-list", m.Cache(h.GetList)).Methods("PUT")
	file.HandleFunc("/get-list", c.GetList).Methods("PUT")
	file.HandleFunc("/get-total-size", c.GetSize).Methods("GET")
	file.HandleFunc("/delete", c.MoveToTrash).Methods("PUT")
	file.HandleFunc("/move", c.MoveFile).Methods("PUT")
	file.HandleFunc("/createFolder", c.CreateFolder).Methods("POST")
	file.HandleFunc("/deleteFolder", c.DeleteFolder).Methods("DELETE")
	file.HandleFunc("/star", c.StarFileOrFolder).Methods("POST")
	file.HandleFunc("/unstar", c.UnStar).Methods("DELETE")

	// implementing rate limiting for the following routes
	// file.Use(middlewares.RateLimiter)
	file.HandleFunc("/download/{type}", h.DownloadFile).Methods("PUT")
	file.HandleFunc("/upload/{type}", h.UploadFile).Methods("POST")
	file.HandleFunc("/rename", h.Rename).Methods("PUT")

	// file trash
	file.HandleFunc("/trash-get", c.GetTrashedFiles).Methods("PUT")
	file.HandleFunc("/trash-delete-forever", h.DeleteFile).Methods("PUT")
	file.HandleFunc("/trash-restore", c.RestoreTrash).Methods("PUT")

	return nil
}
