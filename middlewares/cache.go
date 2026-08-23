package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mayukh551/cloudbox/models"
	"github.com/mayukh551/cloudbox/utils"
)

func Cache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/api/v1/file/get-list" {
			next.ServeHTTP(w, r)
			return
		}

		var payload models.FileListPayload
		var data []models.FileList
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		r.Body.Close()

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		json.Unmarshal(bodyBytes, &data)

		userID, err := utils.GetUserID(r)
		if err != nil || userID == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]any{
				"success": false,
				"error":   "Failed to get userID.",
			})
		}

		cacheKey := fmt.Sprintf("files:%s:%s:%s:%d:%d", userID, payload.Category, payload.Path, payload.Page, payload.Limit)
		var cachedData []models.FileList = nil
		if err = utils.GetCacheObject(r.Context(), cacheKey, &cachedData); err != nil {
			cachedData = nil
		}

		if cachedData != nil {
			fmt.Println("CACHE HIT!!")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(map[string]any{
				"success": true,
				"data":    cachedData,
			})
			return
		}

		fmt.Println("CACHE MISS!!")

		next.ServeHTTP(w, r)

	})
}
