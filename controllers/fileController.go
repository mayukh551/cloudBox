package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {

	var data map[string]any

	json.NewDecoder(r.Body).Decode(&data)

	fileKey := data["file"].(string)
	if fileKey == "" {
		http.Error(w, "missing 'file' query parameter", http.StatusBadRequest)
		return
	}

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-southeast-2"),
	)

	if err != nil {
		http.Error(w, "failed to load AWS config", http.StatusInternalServerError)
		return
	}

	svc := s3.NewFromConfig(cfg)

	// Get the specific object from S3
	result, err := svc.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String("project-cloudbox-assets"),
		Key:    aws.String(fileKey),
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

func UploadFile(w http.ResponseWriter, r *http.Request) {

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
		config.WithRegion("ap-southeast-2"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create S3 service client
	svc := s3.NewFromConfig(cfg)

	_, err = svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("project-cloudbox-assets"),
		Key:    aws.String(handler.Filename),
		Body:   file,
	})

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error while uploading your file", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, fmt.Sprintf("%s Uploaded Successfully!", handler.Filename), http.StatusOK)

}
