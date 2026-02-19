package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/mayukh551/cloudbox/db"
	"github.com/mayukh551/cloudbox/models"
	"github.com/mayukh551/cloudbox/utils"
)

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetUserID(r)

	if err != nil {
		respondWithJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := db.GetUserByID(id)

	if err != nil {
		respondWithJSON(w, "Error getting user details", http.StatusInternalServerError)
		return
	}

	user, err = db.GetUserByID(user.ID)

	if err != nil {
		respondWithJSON(w, "Error getting user details", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, user, http.StatusOK)
}

func UpdateUserDetails(w http.ResponseWriter, r *http.Request) {

	var user models.UpdateUser

	id, err := utils.GetUserID(r)

	if err != nil {
		respondWithJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewDecoder(r.Body).Decode(&user)

	err = db.UpdateUser(id, user)

	if err != nil {
		respondWithJSON(w, "Error updating user details", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, "User updated successfully", http.StatusOK)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	id, err := utils.GetUserID(r)

	if err != nil {
		respondWithJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.DeleteUser(id)

	if err != nil {
		respondWithJSON(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, "User deleted successfully", http.StatusOK)
}
