package service

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// handleCacheStats returns cache statistics as JSON
func (s *Service) handleCacheStats(w http.ResponseWriter, r *http.Request) *appError {
	if !s.cache.Enabled() {
		return writeJSON(w, "application/json", map[string]string{
			"status": "disabled",
		})
	}

	stats := s.cache.Stats()
	return writeJSON(w, "application/json", stats)
}

// handleCacheClear clears the entire cache
func (s *Service) handleCacheClear(w http.ResponseWriter, r *http.Request) *appError {
	if !s.cache.Enabled() {
		return appErrorBadRequest(nil, "Cache is disabled")
	}

	s.cache.Clear()

	return writeJSON(w, "application/json", map[string]string{
		"status":  "ok",
		"message": "Cache cleared",
	})
}

// handleCacheClearLayer clears all tiles for a specific layer
func (s *Service) handleCacheClearLayer(w http.ResponseWriter, r *http.Request) *appError {
	if !s.cache.Enabled() {
		return appErrorBadRequest(nil, "Cache is disabled")
	}

	vars := mux.Vars(r)
	layer := vars["layer"]

	removed := s.cache.ClearLayer(layer)

	return writeJSON(w, "application/json", map[string]interface{}{
		"status":  "ok",
		"message": fmt.Sprintf("Cleared %d tiles for layer %s", removed, layer),
		"removed": removed,
		"layer":   layer,
	})
}
