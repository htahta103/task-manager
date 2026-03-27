package main

import (
	"net/http"
	"strings"
)

func (s *apiServer) withMiddleware(next http.Handler) http.Handler {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure default JSON content type even for 404/405 responses.
		w.Header().Set("Content-Type", "application/json")

		// CORS (very small v1: allow a single origin if configured)
		if s.corsOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", s.corsOrigin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Optional auth: if AUTH_TOKEN is set, require "Authorization: Bearer <token>"
		if s.authToken != "" && strings.HasPrefix(r.URL.Path, "/api/") {
			auth := strings.TrimSpace(r.Header.Get("Authorization"))
			if auth != "Bearer "+s.authToken {
				writeError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
		}

		next.ServeHTTP(w, r)
	})
	return h
}
