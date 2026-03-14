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
		respondWithError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateStruct(data); err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := utils.GetRequestUser(r)

	sharedTo, err := db.GetUserByEmail(data.Email)

	if err != nil {
		respondWithError(w, "User not found!", 404)
		return
	}

	err = db.CreateShare(models.ShareUser{
		SharedTo: sharedTo.ID,
		SharedBy: user.ID,
		FileID:   data.FileID,
	})

	if err != nil {
		respondWithError(w, err.Error(), 500)
		return
	}

	respondWithJSON(w, "", 200)
}

func ListShares(w http.ResponseWriter, r *http.Request) {

	user := utils.GetRequestUser(r)

	shares := db.ListShares(user.ID)

	if shares == nil {
		respondWithError(w, "No shares found!", 404)
		return
	}

	respondWithJSON(w, shares, 200)
}

func ListSharedWithMe(w http.ResponseWriter, r *http.Request) {

	user := utils.GetRequestUser(r)

	shares := db.ListSharedWithMe(user.ID)

	if shares == nil {
		respondWithError(w, "No shares found!", 404)
		return
	}

	respondWithJSON(w, shares, 200)
}
