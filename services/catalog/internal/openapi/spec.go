package openapi

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed openapi.yaml
var specFS embed.FS

func Handler() http.HandlerFunc {
	spec, _ := fs.ReadFile(specFS, "openapi.yaml")
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		w.Header().Set("Content-Disposition", `inline; filename="openapi.yaml"`)
		w.WriteHeader(http.StatusOK)
		w.Write(spec)
	}
}
