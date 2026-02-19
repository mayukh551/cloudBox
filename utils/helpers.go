package utils

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mayukh551/cloudbox/middlewares"
	"github.com/mayukh551/cloudbox/models"
	"golang.org/x/crypto/bcrypt"
)

type JWTClaim struct {
	ID    string
	Email string
	jwt.RegisteredClaims
}

var jwtKey = []byte("secret_key")

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func ValidatePassword(textPassword string, hashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(textPassword))
	return err == nil
}

func GenerateJWTToken(user models.User) (string, error) {

	fmt.Println("in GenerateJWTToken", user)

	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims with user data
	claims := JWTClaim{
		ID:    user.ID,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "your-api",
			Subject:   user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	fmt.Println(tokenString, claims.ID, claims.Email)

	return tokenString, nil

}

func VerifyJWTToken() {

}

func GenerateUUID() string {
	return uuid.New().String()
}

func GetUserID(r *http.Request) (string, error) {
	userCtxt, ok := r.Context().Value("user").(middlewares.RequestUser)
	if !ok {
		return "", errors.New("user not found in req context")
	}
	return userCtxt.ID, nil
}
