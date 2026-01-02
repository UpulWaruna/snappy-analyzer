package main

import (
	"bytes"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"common/logger"

	"github.com/rs/cors"
)

// AnalysisRequest defines the incoming JSON from React
type AnalysisRequest struct {
	URL string `json:"url"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/analyze", analysisHandler)

	// Setup CORS options
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Your React URL
		AllowedMethods: []string{"POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
		Debug:          true, // Useful for troubleshooting
	})

	// Wrap the mux with the CORS middleware
	handler := c.Handler(mux)

	logger.InitSharedLogger("Worker Service")
	slog.Info("Worker Service started on", "port", 8080)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func analysisHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Decode the request
	var req AnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 2. Respond to React immediately
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "processing",
		"message": "Analysis started for " + req.URL,
	})

	// 3. Process in background
	go func(targetURL string) {
		slog.Info("Background analysis: %s", targetURL)

		// A. Render HTML
		renderedHTML, err := GetRenderedHTML(targetURL)
		if err != nil {
			sendErrorToSocket(targetURL, "ChromeDP failed: "+err.Error())
			return
		}

		// B. Parse HTML (This creates our primary Result object)
		parsedResult, err := ParseHTML(bytes.NewReader([]byte(renderedHTML)))
		if err != nil {
			sendErrorToSocket(targetURL, "Parser failed: "+err.Error())
			return
		}

		// C. Fill in the missing URL (since ParseHTML doesn't know it)
		parsedResult.URL = targetURL

		// D. Process Links & Update the parsedResult directly
		checkedLinks := ProcessLinks(targetURL, parsedResult.discoveredLinks)

		for _, l := range checkedLinks {
			if l.IsExternal {
				parsedResult.Links.ExternalCount++
			} else {
				parsedResult.Links.InternalCount++
			}
			if !l.Accessible {
				parsedResult.Links.Inaccessible++
			}
		}

		// 4. Send the populated parsedResult to the socket service
		sendToSocketService(parsedResult)
		slog.Info("Successfully broadcasted results for: %s", targetURL)
	}(req.URL)
}

// Helper for cleaner error reporting
func sendErrorToSocket(targetURL string, message string) {
	sendToSocketService(map[string]interface{}{
		"url":   targetURL,
		"error": map[string]string{"message": message},
	})
}

// sendToSocketService POSTs the result to the broadcaster container
func sendToSocketService(result interface{}) {
	jsonData, err := json.Marshal(result)
	if err != nil {
		slog.Info("Marshal error: %v", err)
		return
	}

	// CHANGE THIS:
	// From: "http://localhost:8081/publish" for local testing
	// To:   "http://socket-service:8081/publish" for Docker networking

	// We use the container name 'socket-service' defined in docker-compose
	socketServiceURL := "http://socket-service:8081/publish"

	resp, err := http.Post(socketServiceURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Info("Failed to reach Socket Service: %v", err)
		return
	}
	defer resp.Body.Close()
}
