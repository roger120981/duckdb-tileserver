package service

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tobilg/duckdb-tileserver/internal/data"
)

// LayersResponse represents the JSON response for the /layers endpoint
type LayersResponse struct {
	Layers []*data.Layer `json:"layers"`
}

// handleLayers returns a list of all available spatial layers
func handleLayers(w http.ResponseWriter, r *http.Request) *appError {
	log.Debug("Layers request")

	// Get catalog instance
	catDB, ok := catalogInstance.(*data.CatalogDB)
	if !ok {
		return appErrorInternal(nil, "Invalid catalog type")
	}

	// Get all layers
	layers, err := catDB.GetLayers()
	if err != nil {
		return appErrorInternal(err, fmt.Sprintf("Error retrieving layers: %v", err))
	}

	response := LayersResponse{
		Layers: layers,
	}

	return writeJSON(w, ContentTypeJSON, response)
}
