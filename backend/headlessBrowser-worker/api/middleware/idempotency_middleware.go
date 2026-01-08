package middleware

import (
	"net/http"
	"sync"
	"time"
)

var (
	idempotencyKeys = make(map[string]time.Time)
	muID            sync.Mutex
)

// add code line for health check by passing through the request if the path is /health
func IdempotencyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		idempotencyKey := r.Header.Get("Idempotency-Key")
		if idempotencyKey == "" {
			http.Error(w, "Missing Idempotency-Key", http.StatusBadRequest)
			return
		}

		muID.Lock()
		defer muID.Unlock()

		if _, exists := idempotencyKeys[idempotencyKey]; exists {
			http.Error(w, "Duplicate request", http.StatusConflict)
			return
		}

		idempotencyKeys[idempotencyKey] = time.Now()

		next.ServeHTTP(w, r)
	})
}
