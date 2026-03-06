package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// UserIDHeader is the header set by Kong after validating the JWT (claim "sub").
// The service trusts this header; Kong is responsible for auth. See docs (e.g. project-plan, Kong JWT plugin).
const UserIDHeader = "X-User-ID"

type contextKey string

const userIDContextKey contextKey = "userID"

// UserIDFromContext returns the authenticated user ID from the request context.
// Set by RequireUserID (which reads the header set by Kong).
func UserIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(userIDContextKey).(string)
	return id
}

// RequireUserID returns a middleware that requires the user ID header (set by Kong after JWT validation)
// and puts it in the request context. Responds 401 if the header is missing or empty.
//
// Kong validates the Bearer JWT and forwards the subject (user ID) in X-User-ID. This service
// does not validate JWT — that is the gateway's responsibility.
//
// For local dev without Kong: call GET /me with header X-User-ID: <user-uuid> (e.g. copy from a registered user).
func RequireUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := strings.TrimSpace(r.Header.Get(UserIDHeader))
		if userID == "" {
			respondUnauthorized(w, "missing user identity (gateway must set "+UserIDHeader+" after JWT validation)")
			return
		}
		ctx := context.WithValue(r.Context(), userIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func respondUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
