package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mayukh551/cloudbox/utils"
)

type File struct {
	Name       string `json:"name"`
	ModifiedAt string `json:"modifiedAt"`
	Size       int64  `json:"size"`
}

type S3Handler struct {
	s3         *s3.Client
	region     string
	bucketName string
}

func NewHandler(s3 *s3.Client) *S3Handler {
	return &S3Handler{
		s3:         s3,
		region:     os.Getenv("S3_REGION"),
		bucketName: os.Getenv("S3_BUCKET_NAME"),
	}
}

func fetchUserID(w http.ResponseWriter, r *http.Request) string {
	userID, err := utils.GetUserID(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}

	return userID
}

func (h *S3Handler) GetList(w http.ResponseWriter, r *http.Request) {

	userID := fetchUserID(w, r)

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(h.region),
	)

	if err != nil {
		http.Error(w, "failed to load AWS config", http.StatusInternalServerError)
		return
	}

	svc := s3.NewFromConfig(cfg)

	result, err := svc.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(h.bucketName),
		Prefix: aws.String(userID),
	})

	if err != nil {
		http.Error(w, "failed to list objects: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var files []File = []File{}
	for _, obj := range result.Contents {
		key := aws.ToString(obj.Key)
		// Extract filename from userID/filename
		parts := strings.Split(key, "/")
		name := key
		if len(parts) > 1 {
			name = parts[1]
		}

		files = append(files, File{
			Name:       name,
			ModifiedAt: obj.LastModified.String(),
			Size:       *obj.Size,
		})
	}

	respondWithJSON(w, files, http.StatusOK)

}

func (h *S3Handler) DownloadFile(w http.ResponseWriter, r *http.Request) {

	userID := fetchUserID(w, r)

	var data map[string]any

	json.NewDecoder(r.Body).Decode(&data)

	fileKey := data["file"].(string)
	if fileKey == "" {
		http.Error(w, "missing 'file' query parameter", http.StatusBadRequest)
		return
	}

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(h.region),
	)

	if err != nil {
		http.Error(w, "failed to load AWS config", http.StatusInternalServerError)
		return
	}

	svc := s3.NewFromConfig(cfg)

	// Get the specific object from S3
	result, err := svc.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(h.bucketName),
		Key:    aws.String(userID + "/" + fileKey),
	})

	if err != nil {
		http.Error(w, "file not found: "+err.Error(), http.StatusNotFound)
		return
	}

	defer result.Body.Close()

	// Set response headers
	filename := path.Base(fileKey) // extract filename from key

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", aws.ToString(result.ContentType))
	if result.ContentLength != nil {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", *result.ContentLength))
	}

	// Stream S3 body directly to the HTTP response
	if _, err := io.Copy(w, result.Body); err != nil {
		log.Println("error streaming file:", err)
	}
}

func (h *S3Handler) UploadFile(w http.ResponseWriter, r *http.Request) {

	userID := fetchUserID(w, r)

	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		http.Error(w, "Error while parsing file", 400)
		return
	}

	file, handler, err := r.FormFile("file")

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}

	defer file.Close()

	// Load AWS config (V2 SDK — no more sessions)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(h.region),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create S3 service client
	svc := s3.NewFromConfig(cfg)

	_, err = svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(h.bucketName),
		Key:    aws.String(userID + "/" + handler.Filename),
		Body:   file,
	})

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error while uploading your file", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, fmt.Sprintf("%s Uploaded Successfully!", handler.Filename), http.StatusOK)

}

func (h *S3Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {

	userID := fetchUserID(w, r)

	var data map[string]any

	json.NewDecoder(r.Body).Decode(&data)

	fileKey := data["file"].(string)
	if fileKey == "" {
		http.Error(w, "missing 'file' query parameter", http.StatusBadRequest)
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(h.region),
	)
	if err != nil {
		log.Fatal(err)
	}

	svc := s3.NewFromConfig(cfg)

	_, err = svc.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(h.bucketName),
		Key:    aws.String(userID + "/" + fileKey),
	})

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error while deleting your file", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, fmt.Sprintf("%s Deleted Successfully!", fileKey), http.StatusOK)
}
