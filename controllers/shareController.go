package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/mayukh551/cloudbox/db"
	"github.com/mayukh551/cloudbox/models"
	"github.com/mayukh551/cloudbox/utils"
)

func Share(w http.ResponseWriter, r *http.Request) {

	var data models.FileShare

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, utils.JSON_DECODE_ERROR, http.StatusBadRequest, err)
		return
	}

	if err := utils.ValidateStruct(data); err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest, err)
		return
	}

	userID := fetchUserID(w, r)

	sharedTo, err := db.GetUserByEmail(data.Email, r.Context())

	if err != nil {
		respondWithError(w, "User not found!", 404, err)
		return
	}

	err = db.CreateShare(models.ShareUser{
		SharedTo: sharedTo.ID,
		SharedBy: userID,
		FileID:   data.FileID,
	}, r.Context())

	if err != nil {
		respondWithError(w, err.Error(), 500, err)
		return
	}

	respondWithJSON(w, nil, 200)
}

// func ListShares(w http.ResponseWriter, r *http.Request) {

// 	userID := fetchUserID(w, r)

// 	shares, err := db.ListShares(userID, r.Context())

// 	if err != nil {
// 		respondWithError(w, "No shares found!", 404, nil)
// 		return
// 	}

// 	respondWithJSON(w, shares, 200)
// }

// func ListSharedWithMe(w http.ResponseWriter, r *http.Request) {

// 	userID := fetchUserID(w, r)

// 	shares := db.ListSharedWithMe(userID, r.Context())

// 	if shares == nil {
// 		respondWithError(w, "No shares found!", 404, nil)
// 		return
// 	}

// 	respondWithJSON(w, shares, 200)
// }

func RemoveSharedFiles(w http.ResponseWriter, r *http.Request) {

	shareID := r.URL.Query().Get("shareID")

	shares := db.RemoveSharedFile(shareID, r.Context())

	if shares == nil {
		respondWithError(w, "No shares found!", 404, nil)
		return
	}

	respondWithJSON(w, shares, 200)
}
