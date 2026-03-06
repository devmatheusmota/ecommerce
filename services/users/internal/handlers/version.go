package handlers

import "os"

// Version is set at build time via ldflags, e.g.:
//
//	go build -ldflags "-X github.com/ecommerce/services/users/internal/handlers.Version=1.0.0"
//
// If not set, it defaults to "dev".
var Version = "dev"

// VersionString returns the version to expose in API responses.
// In dev: set env VERSION (e.g. VERSION=dev-1.0.0) to show a dev version; otherwise uses build-time Version or "dev".
func VersionString() string {
	if v := os.Getenv("VERSION"); v != "" {
		return v
	}
	return Version
}
