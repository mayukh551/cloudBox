package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaim represents the claims in the JWT token
type JWTClaim struct {
	ID    string
	Email string
	jwt.RegisteredClaims
}

type RequestUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

var jwtKey = []byte("secret_key")

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := bearerToken[1]

		claims := &JWTClaim{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				// Validate the signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return jwtKey, nil
			},
		)

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid token signature", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		fmt.Println(claims.ID, claims.Email)

		user := RequestUser{
			ID:    claims.ID,
			Email: claims.Email,
		}

		reqCtxt := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(reqCtxt)

		next.ServeHTTP(w, r)

	})
}
