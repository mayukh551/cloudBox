package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"regexp"
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
		respondWithError(w, err.Error(), http.StatusUnauthorized, err)
		return ""
	}

	return userID
}

func GetList(w http.ResponseWriter, r *http.Request) {

	userID := fetchUserID(w, r)

	var data models.FileListPayload
	json.NewDecoder(r.Body).Decode(&data)

	if data.Search != "" {

		// Handle folder names as well in the search
		files, err := db.SearchFilesByName(userID, data.Search, data.Page, data.Limit, r.Context())

		if err != nil {
			respondWithJSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, files, http.StatusOK)

	} else if data.Category == "general" {

		files, err := db.ListFiles(userID, data.Search, data.Path, data.Page, data.Limit, r.Context())

		if err != nil {
			respondWithJSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// fix filenames
		for i := range files {
			// extract filename
			filename := files[i].Name

			// get rid of duplicate prefixes
			var copyPrefixRegex = regexp.MustCompile(`^Copy_\d+_`)
			original := copyPrefixRegex.ReplaceAllString(filename, "")

			// update filename
			files[i].Name = original
		}

		respondWithJSON(w, files, http.StatusOK)

	} else if data.Category == "shared" {
		shares, err := db.ListShares(userID, data.Path, data.Page, data.Limit, r.Context())

		if err != nil {
			respondWithError(w, "No shares found!", 404, nil)
			return
		}

		respondWithJSON(w, shares, http.StatusOK)

	} else if data.Category == "favorites" {
		starredFiles, err := db.GetStarredFiles(userID, data.Path, data.Page, data.Limit, r.Context())

		if err != nil {
			respondWithError(w, "Failed to get starred files!", http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, starredFiles, http.StatusOK)

	} else {
		respondWithError(w, "Choose the right category to load files!", http.StatusInternalServerError, nil)
	}

}

func GetSize(w http.ResponseWriter, r *http.Request) {

	userID := fetchUserID(w, r)

	totalSize, err := db.GetTotalSize(userID, r.Context())

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, totalSize, 200)
}

func (h *S3Handler) Rename(w http.ResponseWriter, r *http.Request) {

	userID := fetchUserID(w, r)

	var data models.UpdateFileNamePayload
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, utils.JSON_DECODE_ERROR, http.StatusBadRequest, err)
		return
	}

	// update filename on table
	data.UpdatedAt = utils.GetCurrentTime()
	err := db.UpdateFileName(data, r.Context())

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError, err)
		return
	}

	// update file on s3
	oldFileKey := userID + "/" + data.OldTitle
	newFileKey := userID + "/" + data.Name

	fmt.Println(oldFileKey, newFileKey)

	// Copy
	if err := utils.CopyObject(h.s3, h.bucketName, oldFileKey, newFileKey); err != nil {
		respondWithError(w, "Error while copying file in s3 bucket!", http.StatusInternalServerError, err)
		return
	}

	// Delete old file
	if err := utils.DeleteObject(h.s3, h.bucketName, oldFileKey); err != nil {
		respondWithError(w, "Failed to delete files from cloud!", http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, "File renamed successfully", http.StatusOK)

}

func (h *S3Handler) DownloadFile(w http.ResponseWriter, r *http.Request) {

	userID := fetchUserID(w, r)

	var data map[string]any

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, utils.JSON_DECODE_ERROR, http.StatusBadRequest, err)
		return
	}

	fileKey := data["file"].(string)
	if fileKey == "" {
		respondWithError(w, "missing 'file' query parameter", http.StatusBadRequest, nil)
		return
	}

	fileKey = userID + "/" + fileKey

	// create presign client to get presign URL for download
	url, err := utils.PresignGetObject(h.s3, h.bucketName, fileKey)

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, models.PreSignedResponse{
		Key: fileKey,
		Url: url,
	}, http.StatusOK)

}

// creates a new record in files Table separately without any upload
func CreateFolder(w http.ResponseWriter, r *http.Request) {

	userID := fetchUserID(w, r)

	var data models.CreateFolderPayload
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, utils.JSON_DECODE_ERROR, http.StatusInternalServerError, nil)
		return
	}

	folderName := data.Name
	path := data.Path

	err := db.CreateFile(
		models.CreateFile{
			ID:     utils.GenerateUUID(),
			Name:   folderName,
			Path:   path,
			Type:   "Folder",
			Size:   0, // size will be updated later via a separate endpoint after new upload is complete
			UserID: userID,
		},
		r.Context(),
	)

	if err != nil {
		respondWithError(w, "Error while creating folder record in database", http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, "Folder created successfully", http.StatusOK)
}

// DeleteFolder takes folder name from req query
// searches the path column that contains this folder name as substring
// if found it deletes it

// Example:
// Folder1
// ├── Folder2
// │   ├── fileA.txt
// │   └── SubFolder1
// │       └── fileB.txt
// └── Folder3
//     └── fileC.txt

// if user requests to delte Folder2, files like fileA.txt, SubFolder1, fileB.txt will be deleted as they come under Folder2
func DeleteFolder(w http.ResponseWriter, r *http.Request) {

	userID := fetchUserID(w, r)

	folderName := r.URL.Query().Get("folderName")
	id := r.URL.Query().Get("id") // ID of the folder

	if err := db.DeleteFolder(folderName, id, userID, r.Context()); err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, "Files and Folder containing it is successfully deleted!", http.StatusOK)

}

// updates the path of a file/folder in the database
func MoveFile(w http.ResponseWriter, r *http.Request) {

	var data models.MoveFilePayload
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, utils.JSON_DECODE_ERROR, http.StatusBadRequest, err)
		return
	}

	data.UpdatedAt = utils.GetCurrentTime()

	if err := db.UpdatePath(data, r.Context()); err != nil {
		respondWithError(w, "Error while updating file path", http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, "File moved successfully", http.StatusOK)
}

// creates a presign URL for upload and sends it back to frontend.
// For already existing file, it updates the path first (This only happens if an empty folder existed with no assets inside)
// and then creates a presign URL for new upload and send back to the frontend
func (h *S3Handler) UploadFile(w http.ResponseWriter, r *http.Request) {

	userID := fetchUserID(w, r)

	var presignPayload models.PreSignedBody

	if err := json.NewDecoder(r.Body).Decode(&presignPayload); err != nil {
		respondWithError(w, utils.JSON_DECODE_ERROR, http.StatusInternalServerError, nil)
		return
	}

	// deconstruct fields from struct
	path := presignPayload.Path
	filename := presignPayload.Filename
	fileID := presignPayload.FileID
	size := presignPayload.Size

	totalSize, err := db.GetTotalSize(userID, r.Context())
	if err != nil {
		respondWithError(w, "Error while getting total storage size", http.StatusInternalServerError, err)
		return
	}

	const maxAllowedStorage = 1024 * 1024 * 1024 // 1GB
	if (totalSize + size) > maxAllowedStorage {
		respondWithError(w, "Storage limit exceeded: total usage cannot exceed 1GB", http.StatusBadRequest, nil)
		return
	}

	// if FileID is provided, then we are updating path and filename
	if presignPayload.FileID != "" {
		if err := db.UpdateFile(
			models.CreateFile{
				ID:        fileID,
				Name:      filename,
				Path:      path,
				UpdatedAt: fmt.Sprintf("%d", time.Now()),
				Size:      size,
				Type:      "File",
			},
			r.Context(),
		); err != nil {
			respondWithError(w, "Error while updating file path", http.StatusInternalServerError, err)
			return
		}
	}

	ifFileExists, err := db.CheckIfFileNameExists(presignPayload.Filename, path, r.Context())

	if err != nil {
		respondWithError(w, "Failed to check if file exists", http.StatusInternalServerError, err)
		return
	}

	// if filename with the same already exists, then return error, two or more files with same name cannot be stored in the same path
	if ifFileExists["samePathExists"] == true {
		respondWithError(w, fmt.Sprintf("Another filename with %s already exists.", filename), http.StatusBadRequest, nil)
		return
	}

	// to handle duplicate names in the same path for multiple files in s3
	if ifFileExists["diffPathExists"] == true {
		filename = fmt.Sprintf("Copy_%d_%s", rand.IntN(9000)+1000, filename)
	}

	key := userID + "/" + filename

	url, err := utils.PresignPutObject(h.s3, h.bucketName, key, presignPayload.ContentType)

	if err != nil {
		respondWithError(w, "Error while creating a presigned URL", http.StatusInternalServerError, err)
		return
	}

	if err = db.CreateFile(
		models.CreateFile{
			ID:     utils.GenerateUUID(),
			Name:   filename,
			Path:   path,
			Type:   "file",
			Size:   size, // size will be updated later via a separate endpoint after upload is complete
			UserID: userID,
		},
		r.Context(),
	); err != nil {
		respondWithError(w, "Error while creating file record in database", http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, models.PreSignedResponse{
		Key: key,
		Url: url,
	}, http.StatusOK)

}

// Update a single file record with updated fields and must have FileID
func UpdateFile(w http.ResponseWriter, r *http.Request) {

	var updatedData models.CreateFile

	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		respondWithError(w, utils.JSON_DECODE_ERROR, http.StatusInternalServerError, nil)
		return
	}

	if err := db.UpdateFile(updatedData, r.Context()); err != nil {
		respondWithError(w, "Error while updating file metadata", http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, "File Updated Successfully", http.StatusOK)

}

func MoveToTrash(w http.ResponseWriter, r *http.Request) {

	var data models.TrashPayload

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		respondWithError(w, utils.JSON_DECODE_ERROR, http.StatusBadRequest, err)
		return
	}

	if len(data.Files) == 0 {
		respondWithError(w, "no files provided", http.StatusBadRequest, nil)
		return
	}

	if err := db.SetIsTrash(data.Files, true, r.Context()); err != nil {
		respondWithError(w, "Error moving files to trash", http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, "Files moved to trash", http.StatusOK)

}

// Delete File from files table and from s3 bucket
func (h *S3Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {

	userID := fetchUserID(w, r)

	var data []models.DeleteFilePayload

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		respondWithError(w, utils.JSON_DECODE_ERROR, http.StatusBadRequest, err)
		return
	}

	files := data
	if len(files) == 0 {
		respondWithError(w, "'files' field cannot be empty", http.StatusBadRequest, err)
		return
	}

	// *************************

	var errorFileNames []string

	var wg sync.WaitGroup

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
				errorFileNames = append(errorFileNames, file.Key)
			}

			// delete from Files Table
			if err := db.DeleteFile(file.Id, r.Context()); err != nil {
				respondWithError(w, "Failed to delete file from db", http.StatusInternalServerError, err)
			}
		}()
	}

	wg.Wait()

	if len(errorFileNames) > 0 {
		respondWithError(w, errorFileNames, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, fmt.Sprintf("File(s) Deleted Successfully!"), http.StatusOK)
}

func StarFileOrFolder(w http.ResponseWriter, r *http.Request) {

	// type
	// name

	// if type file, then go
	// if type folder, then folder name

	var data models.StarFilePayload

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, utils.JSON_DECODE_ERROR, http.StatusInternalServerError, err)
		return
	}

	userID := fetchUserID(w, r)
	var wg sync.WaitGroup

	var err error

	for i := 0; i < len(data.IDs); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = db.AddStar(data.IDs[i], userID, r.Context())
		}()
	}

	if err != nil {
		respondWithError(w, "Error while saving file as favourite.", 500, err)
		return
	}

	wg.Wait()

	respondWithJSON(w, "File marked as favourite.", 200)

}

func UnStar(w http.ResponseWriter, r *http.Request) {

	fileID := r.URL.Query().Get("fileID")
	userID := fetchUserID(w, r)

	if err := db.RemoveStar(fileID, userID, r.Context()); err != nil {
		respondWithError(w, "Error while saving file as favourite.", 500, err)
		return
	}

	respondWithJSON(w, "File removed from favourite.", 200)

}

// Handling trashes
func GetTrashedFiles(w http.ResponseWriter, r *http.Request) {

	userID := fetchUserID(w, r)

	fmt.Println("yooo")

	var data models.FileListPayload

	err := json.NewDecoder(r.Body).Decode(&data)

	trashFiles, err := db.GetTrashFiles(userID, data.Page, data.Limit, r.Context())

	if err != nil {
		respondWithError(w, "Error fetching trash files!!", http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, trashFiles, http.StatusOK)

}

func RestoreTrash(w http.ResponseWriter, r *http.Request) {

	ids := r.URL.Query()["fileID"]

	// Restore: set isTrash = false
	if err := db.SetIsTrash(ids, false, r.Context()); err != nil {
		respondWithError(w, "Error restoring files from trash", http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, "Files restored", http.StatusOK)
}
