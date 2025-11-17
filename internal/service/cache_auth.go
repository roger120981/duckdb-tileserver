package service

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tobilg/duckdb-tileserver/internal/conf"
)

const (
	headerAPIKey = "X-API-Key"
)

// cacheAuthMiddleware validates API key for cache endpoints
func cacheAuthMiddleware(next appHandler) appHandler {
	return func(w http.ResponseWriter, r *http.Request) *appError {
		// Get configured API key
		configuredKey := conf.Configuration.Cache.ApiKey

		// If no API key is configured, allow access (public mode)
		if configuredKey == "" {
			log.Debug("Cache endpoint accessed without authentication (public mode)")
			return next(w, r)
		}

		// API key is configured, validate the request
		providedKey := r.Header.Get(headerAPIKey)

		// Check if key was provided
		if providedKey == "" {
			log.Warnf("Cache endpoint accessed without API key from %s", r.RemoteAddr)
			return appErrorUnauthorized(nil, "API key required. Provide X-API-Key header.")
		}

		// Validate the key
		if providedKey != configuredKey {
			log.Warnf("Cache endpoint accessed with invalid API key from %s", r.RemoteAddr)
			return appErrorForbidden(nil, "Invalid API key")
		}

		// Authentication successful
		log.Debugf("Cache endpoint accessed with valid API key from %s", r.RemoteAddr)
		return next(w, r)
	}
}
