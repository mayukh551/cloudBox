package routers

import (
	"github.com/gorilla/mux"
	c "github.com/mayukh551/cloudbox/controllers"
	m "github.com/mayukh551/cloudbox/middlewares"
	"github.com/mayukh551/cloudbox/utils"
)

// fileRoutes registers every /api/v1/file/* endpoint. It returns an error
// because it needs to load AWS config to build the S3-backed handler (h);
// see the caller-side note in router.go about that error currently being
// ignored.
func fileRoutes(api *mux.Router) error {
	// File APIs
	file := api.PathPrefix("/file").Subrouter()

	// All routes under /file require a valid auth token.
	file.Use(m.Authenticate)

	// Load AWS credentials/config once at startup and reuse the S3 client
	// across requests (avoids reconnecting per-request).
	s3, err := utils.LoadAWSConfig()
	if err != nil {
		return err
	}

	h := c.NewHandler(s3)

	// --- Plain DB-backed routes (no S3 interaction needed) ---

	// file.Use(m.Cache)
	file.HandleFunc("/get-list", c.GetList).Methods("PUT")
	file.HandleFunc("/get-total-size", c.GetSize).Methods("GET")
	file.HandleFunc("/delete", c.MoveToTrash).Methods("PUT")
	file.HandleFunc("/move", c.MoveFile).Methods("PUT")
	file.HandleFunc("/createFolder", c.CreateFolder).Methods("POST")
	file.HandleFunc("/deleteFolder", c.DeleteFolder).Methods("DELETE")
	file.HandleFunc("/star", c.StarFileOrFolder).Methods("POST")
	file.HandleFunc("/unstar", c.UnStar).Methods("DELETE")

	// --- S3-backed routes (use the handler bound to the S3 client) ---
	// implementing rate limiting for the following routes
	file.HandleFunc("/download/{type}", h.DownloadFile).Methods("PUT")
	file.HandleFunc("/upload/{type}", h.UploadFile).Methods("POST")
	file.HandleFunc("/rename", h.Rename).Methods("PUT")

	// --- Trash management ---
	file.HandleFunc("/trash-get", c.GetTrashedFiles).Methods("PUT")
	file.HandleFunc("/trash-delete-forever", h.DeleteFile).Methods("PUT")
	file.HandleFunc("/trash-restore", c.RestoreTrash).Methods("PUT")

	return nil
}
