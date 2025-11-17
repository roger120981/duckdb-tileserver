package service

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tobilg/duckdb-tileserver/internal/ui"
)

// serveMapViewer serves the HTML map viewer page
func serveMapViewer(w http.ResponseWriter, r *http.Request) *appError {
	log.Debug("Map viewer request")

	// Load the template
	templ, err := ui.LoadTemplate("index.gohtml")
	if err != nil {
		return appErrorInternal(err, "Error loading map viewer template")
	}

	// Execute the standalone template directly (it's a complete HTML page)
	w.Header().Set("Content-Type", ContentTypeHTML)
	w.WriteHeader(http.StatusOK)
	err = templ.Execute(w, nil)
	if err != nil {
		return appErrorInternal(err, "Error rendering map viewer")
	}

	return nil
}
