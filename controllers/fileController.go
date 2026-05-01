package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mayukh551/cloudbox/db"
	"github.com/mayukh551/cloudbox/models"
	"github.com/mayukh551/cloudbox/utils"
)

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

	if err != nil || userID == "" {
		respondWithError(w, err.Error(), http.StatusUnauthorized)
		return ""
	}

	return userID
}

func (h *S3Handler) GetList(w http.ResponseWriter, r *http.Request) {

	userID := fetchUserID(w, r)

	files, err := db.ListFiles(userID, r.Context())

	fmt.Println(files)

	if err != nil {
		respondWithError(w, "failed to list objects", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, files, http.StatusOK)

}

func (h *S3Handler) Rename(w http.ResponseWriter, r *http.Request) {

	userID := fetchUserID(w, r)

	var data models.UpdateFileNamePayload
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, utils.JSON_DECODE_ERROR, http.StatusBadRequest)
		return
	}

	// update filename on table
	data.UpdatedAt = time.Now().Format(time.RFC3339) // TODO: need to recheck this part
	err := db.UpdateFileName(data, r.Context())

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
	}

	// update file on s3
	oldFileKey := userID + "/" + data.OldTitle

	// Copy
	_, err = h.s3.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     aws.String(h.bucketName),
		CopySource: aws.String(oldFileKey),
		Key:        aws.String(data.Title),
	})

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
	}

	// Delete old file
	if err = utils.DeleteObject(h.s3, h.bucketName, oldFileKey); err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
	}

	respondWithJSON(w, "File renamed successfully", http.StatusOK)

}

func (h *S3Handler) DownloadFile(w http.ResponseWriter, r *http.Request) {

	userID := fetchUserID(w, r)

	var data map[string]any

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, utils.JSON_DECODE_ERROR, http.StatusBadRequest)
		return
	}

	fileKey := data["file"].(string)
	if fileKey == "" {
		respondWithError(w, "missing 'file' query parameter", http.StatusBadRequest)
		return
	}

	fileKey = userID + "/" + fileKey

	// create presign client to get presign URL for download
	url, err := utils.PresignGetObject(h.s3, h.bucketName, fileKey)

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, models.PreSignedResponse{
		Key: fileKey,
		Url: url,
	}, http.StatusOK)

}

func (h *S3Handler) UploadFile(w http.ResponseWriter, r *http.Request) {

	userID := fetchUserID(w, r)

	var presignPayload models.PreSignedBody
	json.NewDecoder(r.Body).Decode(&presignPayload)

	key := userID + "/" + presignPayload.Filename

	url, err := utils.PresignPutObject(h.s3, h.bucketName, key, presignPayload.ContentType)

	if err != nil {
		respondWithError(w, "Error while creating a presigned URL", http.StatusInternalServerError)
		return
	}

	err = db.CreateFile(
		models.CreateFile{
			ID:     utils.GenerateUUID(),
			Title:  presignPayload.Filename,
			Type:   "file",
			Size:   presignPayload.Size, // size will be updated later via a separate endpoint after upload is complete
			UserID: userID,
		},
		r.Context(),
	)

	if err != nil {
		respondWithError(w, "Error while creating file record in database", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, models.PreSignedResponse{
		Key: key,
		Url: url,
	}, http.StatusOK)

}

func (h *S3Handler) UpdateFile(w http.ResponseWriter, r *http.Request) {

	var updatedData models.CreateFile

	// TODO: field validation
	json.NewDecoder(r.Body).Decode(&updatedData)

	if err := db.UpdateFile(updatedData.ID, updatedData, r.Context()); err != nil {
		respondWithError(w, "Error while updating file metadata", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, "File Updated Successfully", http.StatusOK)

}

func (h *S3Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {

	userID := fetchUserID(w, r)

	var data []models.DeleteFilePayload

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		respondWithError(w, "Error while decoding JSON", http.StatusBadRequest)
		return
	}

	files := data
	if len(files) == 0 {
		respondWithError(w, "'files' field cannot be empty", http.StatusBadRequest)
		return
	}

	// *************************

	var errorFileNames []string

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, file := range files {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// delete from s3 bucket
			_, err = h.s3.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
				Bucket: aws.String(h.bucketName),
				Key:    aws.String(userID + "/" + file.Key),
			})

			if err != nil {
				mu.Lock()
				errorFileNames = append(errorFileNames, file.Key)
				mu.Unlock()
			}

			// delete from Files Table
			if err := db.DeleteFile(file.Id, r.Context()); err != nil {
				respondWithError(w, "Failed to delete file from db", http.StatusInternalServerError)
			}
		}()
	}

	wg.Wait()

	if len(errorFileNames) > 0 {
		respondWithError(w, errorFileNames, http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, fmt.Sprintf("File(s) Deleted Successfully!"), http.StatusOK)
}
