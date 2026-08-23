package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mayukh551/cloudbox/utils"
	"github.com/redis/go-redis/v9"
)

const (
	CapTokens    = 10
	RefillPerSec = 5.0 // tokens added per second
	ReqCost      = 1
)

// Atomic refill+consume in one round trip.
// KEYS[1] = bucket key
// ARGV[1] = capacity, ARGV[2] = refill rate/sec, ARGV[3] = cost, ARGV[4] = now (unix seconds, float)
var rateLimitScript = redis.NewScript(`
local key = KEYS[1]
local cap = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local cost = tonumber(ARGV[3])
local now = tonumber(ARGV[4])

local data = redis.call("HMGET", key, "tokens", "ts")
local tokens = tonumber(data[1])
local ts = tonumber(data[2])

if tokens == nil then
  tokens = cap
  ts = now
end

local elapsed = math.max(0, now - ts)
tokens = math.min(cap, tokens + elapsed * rate)

local allowed = 0
if tokens >= cost then
  tokens = tokens - cost
  allowed = 1
end

redis.call("HMSET", key, "tokens", tokens, "ts", now)
redis.call("EXPIRE", key, math.ceil(cap / rate) + 1)

return allowed
`)

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := utils.GetUserID(r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]any{"success": false, "data": "Error extracting userID"})
			return
		}

		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}
		if idx := strings.LastIndex(ip, ":"); idx != -1 {
			ip = ip[:idx]
		}

		identity := userID
		if identity == "" {
			identity = ip // fallback for unauthenticated requests
		}

		cacheKey := fmt.Sprintf("rate_limit:%s:%s", identity, r.URL.Path)
		now := float64(time.Now().UnixNano()) / 1e9

		allowed, err := rateLimitScript.Run(r.Context(), utils.RDB, []string{cacheKey},
			CapTokens, RefillPerSec, ReqCost, now).Int()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{"success": false, "data": "Error checking rate limit"})
			return
		}

		if allowed == 1 {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]any{"success": false, "data": "Too many attempts, try later!"})
	})
}
