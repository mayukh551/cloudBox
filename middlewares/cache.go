package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mayukh551/cloudbox/models"
	"github.com/mayukh551/cloudbox/utils"
)

func Cache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var data models.FileListPayload
		json.NewDecoder(r.Body).Decode(&data)

		userID, err := utils.GetUserID(r)
		if err != nil || userID == "" {
			json.NewEncoder(w).Encode(map[string]any{
				"success": false,
				"error":   "Failed to get userID.",
			})
		}

		cacheKey := fmt.Sprintf("files:%s:%s:%d:%d", &userID)

		utils.GetCacheObject(r.Context(), cacheKey)

	})
}
