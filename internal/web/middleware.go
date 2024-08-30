package web

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/germandv/apio/internal/tokenauth"
)

var (
	trueClientIP  = http.CanonicalHeaderKey("True-Client-IP")
	xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
	xRealIP       = http.CanonicalHeaderKey("X-Real-IP")
)

// middleware applies common middleware to all routes.
func middleware(next http.Handler, logger *slog.Logger, auth tokenauth.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Panic-recovery middleware.
		// Recovers panics in other middlewares. It is not needed in API handlers as `ApiFunc` already calls `recover()`.
		defer func() {
			err := recover()
			if err != nil {
				logger.Warn("panic recovered in middleware", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()

		// Add real IP to request.
		var ip string
		if tcip := r.Header.Get(trueClientIP); tcip != "" {
			ip = tcip
		} else if xrip := r.Header.Get(xRealIP); xrip != "" {
			ip = xrip
		} else if xff := r.Header.Get(xForwardedFor); xff != "" {
			i := strings.Index(xff, ", ")
			if i == -1 {
				i = len(xff)
			}
			ip = xff[:i] // first IP in the comma-separated list
		}
		if ip != "" {
			r.RemoteAddr = ip
		}

		// Add security headers.
		// Prevent sensitive information from being cached.
		w.Header().Set("Cache-Control", "no-store")
		// To protect against drag-and-drop style clickjacking attacks.
		w.Header().Set("Content-Security-Policy", "frame-ancestors 'none'")
		w.Header().Set("X-Frame-Options", "DENY")
		// To prevent browsers from performing MIME sniffing, and inappropriately interpreting responses as HTML.
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// Require connections over HTTPS and to protect against spoofed certificates.
		w.Header().Set("Strict-Transport-Security", "max-age=31536000")

		// Enable CORS.
		// Add the "Vary" response header to let caches know that the response
		// will vary according to the request headers `Origin`, `Access-Control-Request-Method` and `Authorization`.
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Authorization")

		// Allow all origins
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// If method is OPTIONS and the header "Access-Control-Request-Method"
		// is present, treat as a preflight request.
		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			// No need to include CORS-safe methods: HEAD, GET and POST, but we do it for explicitness
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE, HEAD, GET, POST")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

			// Cache preflight responses for 10 minutes to save some requests.
			w.Header().Set("Access-Control-Max-Age", "600")

			// Since it's a preflight, send 200 here and don't execute next handlers.
			w.WriteHeader(http.StatusOK)
			return
		}

		// Look for JWT in Authorization header.
		// If it exists, validate it and save claims to request context.
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			if !strings.HasPrefix(authHeader, "Bearer ") {
				w.Header().Set("WWW-Authenticate", "Bearer")
				unauthorized(w, r, logger)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := auth.Validate(token)
			if err != nil {
				unauthorized(w, r, logger)
				return
			}

			userID := claims["sub"].(string)
			role := claims["role"].(string)
			u := CtxUser{ID: userID, Role: role}
			r = r.WithContext(context.WithValue(r.Context(), ctxUserKey, u))
		}

		next.ServeHTTP(w, r)
	})
}

// Helper to send a 401 response.
func unauthorized(w http.ResponseWriter, r *http.Request, logger *slog.Logger) {
	WriteJSON(w, Envelope{"error": "missing, invalid or expired JWT"}, http.StatusUnauthorized)
	logger.Info(
		"API handler finished",
		"method", r.Method,
		"path", r.URL.Path,
		"status", http.StatusUnauthorized,
		"took_ms", 0,
		"ip", r.RemoteAddr,
		"err", "missing or invalid JWT",
	)
}

// RequireUser is a middleware that ensures the user is authenticated.
func RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := GetUser(r.Context())
		if err != nil || u.ID == "" {
			WriteJSON(w, Envelope{"error": "requires authenticated user"}, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireAdmin is a middleware that ensures the user is authenticated and has admin role.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := GetUser(r.Context())
		if err != nil || u.ID == "" || u.Role != "admin" {
			WriteJSON(w, Envelope{"error": "admin only"}, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
