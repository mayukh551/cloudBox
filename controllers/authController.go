package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mayukh551/cloudbox/db"
	"github.com/mayukh551/cloudbox/models"
	"github.com/mayukh551/cloudbox/utils"
)

func SignUp(w http.ResponseWriter, r *http.Request) {

	var data models.CreateUser
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateStruct(data); err != nil {
		fmt.Println(err)
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByEmail(data.Email, r.Context())

	if user != nil {
		respondWithJSON(w, "User already exists!", http.StatusBadRequest)
		return
	}

	// hash password
	hash, err := utils.HashPassword(data.Password)

	if err != nil {
		respondWithJSON(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	data.Password = hash

	data.ID = utils.GenerateUUID()

	if data.ID == "" {
		respondWithJSON(w, "Error generating UUID", http.StatusInternalServerError)
		return
	}

	user, err = db.CreateUser(data, r.Context())

	if err != nil {
		fmt.Println(err)
		respondWithJSON(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	token, err := utils.GenerateJWTToken(*user)

	if err != nil {
		respondWithJSON(w, "Error while generating token.", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, map[string]any{
		"token": token,
	}, 200)
}

func Login(w http.ResponseWriter, r *http.Request) {

	var data models.LoginUser

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateStruct(data); err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := db.VerifyUser(data.Email, data.Password, r.Context())

	if err != nil {
		respondWithJSON(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWTToken(user)

	if err != nil {
		respondWithJSON(w, "Error while generating token.", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, map[string]any{
		"token": token,
	}, 200)
}
