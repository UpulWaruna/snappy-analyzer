package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// LoggingMiddleware logs the incoming HTTP requests and their body
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Read the request body
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
		}

		// Restore the request body to its original state
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Log the request details
		//log.Printf("%s %s %v %s", r.Method, r.URL.Path, time.Since(start), string(bodyBytes))
		slog.Info("incoming request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
			"body", string(bodyBytes),
		)

		next.ServeHTTP(w, r)
	})
}
