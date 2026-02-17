package middlewares

import (
	"assetra/security"
	"assetra/statesvc/restutil"
	"log"
	"net/http"
	"time"
)

func LogRequests(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		next(w, r)
		// Don't log sensitive information like query params or headers
		log.Printf(`{"method": "%s", "path": "%s", "duration": "%v", "status": "completed"}`,
			r.Method, r.URL.Path, time.Since(t))
	}
}

// SecurityHeaders adds essential security headers to all responses
func SecurityHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")
		
		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		
		// XSS Protection (legacy but still useful)
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		
		// Referrer Policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Content Security Policy - strict policy
		w.Header().Set("Content-Security-Policy", 
			"default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'")
		
		// Permissions Policy (formerly Feature Policy)
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		
		// Remove server identification
		w.Header().Del("Server")
		w.Header().Del("X-Powered-By")
		
		next(w, r)
	}
}

// Authenticate validates JWT token and ensures it's not expired
func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := security.ExtractToken(r)
		if err != nil {
			log.Printf("[SECURITY] Authentication failed - missing or invalid token format from IP: %s", r.RemoteAddr)
			restutil.WriteError(w, http.StatusUnauthorized, restutil.ErrUnauthorized)
			return
		}
		
		// Parse and validate token
		token, err := security.ParseToken(tokenString)
		if err != nil {
			log.Printf("[SECURITY] Authentication failed - token parsing error from IP: %s", r.RemoteAddr)
			restutil.WriteError(w, http.StatusUnauthorized, restutil.ErrUnauthorized)
			return
		}
		
		if !token.Valid {
			log.Printf("[SECURITY] Authentication failed - invalid token from IP: %s", r.RemoteAddr)
			restutil.WriteError(w, http.StatusUnauthorized, restutil.ErrUnauthorized)
			return
		}
		
		// Additional validation via payload
		payload, err := security.NewTokenPayload(tokenString)
		if err != nil {
			log.Printf("[SECURITY] %s, Authentication failed - invalid token payload from IP: %s",err, r.RemoteAddr)
			restutil.WriteError(w, http.StatusUnauthorized, restutil.ErrUnauthorized)
			return
		}
		
		// Check if token is expired (defense in depth)
		if time.Now().After(payload.ExpiresAt) {
			log.Printf("[SECURITY] Authentication failed - expired token for user: %s", payload.UserId)
			restutil.WriteError(w, http.StatusUnauthorized, restutil.ErrUnauthorized)
			return
		}

		next(w, r)
	}
}
