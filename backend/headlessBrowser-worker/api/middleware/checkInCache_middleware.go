package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"headlessBrowser-worker/adapter/external"
	"headlessBrowser-worker/domain/service"
	"io"
	"log/slog"
	"net/http"
)

func CheckInCacheMiddleware(dataService *service.AnalysisDataService, l *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Only apply cache to the analyze endpoint (avoid touching /health etc.)
			if r.URL.Path != "/analyze" || r.Method != http.MethodPost {
				next.ServeHTTP(w, r)
				return
			}

			// 1. Get URL from request (assuming POST form or JSON)
			// Note: For HTMX, it's usually r.FormValue("url")
			targetURL := extractURL(r)
			if targetURL == "" {
				next.ServeHTTP(w, r)
				return
			}
			l.Info("from Cache middleware")

			// 2. Check Service/Repo
			cachedResult, err := dataService.RetrieveAnalysisResult(targetURL)
			l.Info("Cache lookup for URL:", "url", targetURL, "found", cachedResult != nil)
			if err == nil && cachedResult != nil {
				// 3. CACHE HIT: Return the result directly!
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Cache", "HIT") // Good for debugging

				// If using HTMX, you'd render the template here instead of JSON
				json.NewEncoder(w).Encode(map[string]interface{}{
					"status": "success",
					"data":   cachedResult,
					"source": "cache",
				})
				//call the analysys publishmethod here to notify the client about the cached result

				publisher := &external.SocketAdapter{Endpoint: "http://socket-service:8081/publish"}
				err = publisher.Publish(cachedResult)
				if err != nil {
					fmt.Println("Error publishing cached result:", err)
					l.Info("Error publishing cached result:", "err", err)
				}
				return // <--- STOP HERE: Don't call the next handler (Analysis)
			}

			// 4. CACHE MISS: Continue to the Analysis Handler
			w.Header().Set("X-Cache", "MISS")
			next.ServeHTTP(w, r)
		})
	}
}

func extractURL(r *http.Request) string {
	// 1) Query param support (optional)
	if u := r.URL.Query().Get("url"); u != "" {
		return u
	}

	// 2) If JSON, read and restore body
	ct := r.Header.Get("Content-Type")
	if len(ct) >= len("application/json") && ct[:len("application/json")] == "application/json" {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return ""
		}
		// Restore body for downstream handler
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var payload struct {
			URL string `json:"url"`
		}
		if err := json.Unmarshal(bodyBytes, &payload); err != nil {
			return ""
		}
		return payload.URL
	}

	// 3) Fallback: form value (works for x-www-form-urlencoded)
	return r.FormValue("url")
}
