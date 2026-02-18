package routes

import (
	"assetra/statesvc/middlewares"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Route struct {
	Method       string
	Path         string
	Handler      http.HandlerFunc
	AuthRequired bool
}

func Install(router *mux.Router, routeList []*Route) {
	for _, route := range routeList {
		handler := route.Handler
		
		// Apply security headers (first)
		handler = middlewares.SecurityHeaders(handler)
		
		// Apply authentication if required
		if route.AuthRequired {
			handler = middlewares.Authenticate(handler)
		}
		
		// Apply logging (last)
		handler = middlewares.LogRequests(handler)
		
		router.HandleFunc(route.Path, handler).Methods(route.Method)
	}
}

func WithCORS(router *mux.Router) http.Handler {
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Accept", "Authorization"})
	
	// Default to localhost only for development
	allowedOrigins := []string{"http://localhost:3000"}
	
	env := os.Getenv("ENV")
	if env == "production" {
		// In production, use specific domain(s) from environment
		prodOrigins := os.Getenv("ALLOWED_ORIGINS")
		if prodOrigins != "" {
			allowedOrigins = []string{prodOrigins}
		} else {
			// Fail secure - no origins allowed if not configured
			allowedOrigins = []string{}
		}
	}
	
	origins := handlers.AllowedOrigins(allowedOrigins)
	methods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions})
	credentials := handlers.AllowCredentials() // Required for httpOnly cookies
	
	return handlers.CORS(headers, origins, methods, credentials)(router)
}
