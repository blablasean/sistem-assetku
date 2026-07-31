package middlewares

import (
	"log"
	"net/http"
	"time"
)

// responseWriterInterceptor captures the status code of the response
type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

func (rwi *responseWriterInterceptor) WriteHeader(code int) {
	rwi.statusCode = code
	rwi.ResponseWriter.WriteHeader(code)
}

// RequestLoggerMiddleware logs HTTP requests in a clean, high-performance, lightweight format
func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Intercept response status code
		rwi := &responseWriterInterceptor{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default status code
		}

		next.ServeHTTP(rwi, r)

		duration := time.Since(start)
		path := r.URL.Path

		// Filter high-frequency background heartbeat requests to prevent console spam
		if path == "/auth/me" || path == "/auth/heartbeat" {
			return
		}

		// Status indicator icon
		icon := "✅"
		if rwi.statusCode >= 400 && rwi.statusCode < 500 {
			icon = "⚠️"
		} else if rwi.statusCode >= 500 {
			icon = "❌"
		}

		log.Printf("%s [%s] %s | Status: %d | Time: %v | Client: %s",
			icon,
			r.Method,
			path,
			rwi.statusCode,
			duration,
			r.RemoteAddr,
		)
	})
}
