package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetFile(w http.ResponseWriter, r *http.Request) {
	// Load AWS config (V2 SDK — no more sessions)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-southeast-2"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create S3 service client
	svc := s3.NewFromConfig(cfg)

	// Get the list of items
	resp, err := svc.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String("project-cloudbox-assets"),
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", item.Size)
		fmt.Println("Storage class:", item.StorageClass)
		fmt.Println("")
	}

	fmt.Println("Found", len(resp.Contents), "items in bucket", "project-cloudbox-assets")
	fmt.Println("")
	respondWithJSON(w, "success", http.StatusOK)
}

func UploadFile(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)

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
