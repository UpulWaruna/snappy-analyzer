package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Define the analysis endpoint
	http.HandleFunc("/analyze", analysisHandler)

	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func analysisHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Enable CORS for React
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Decode user request
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// 3. Fetch the target URL
	resp, err := http.Get(req.URL)
	if err != nil {
		sendError(w, http.StatusServiceUnavailable, fmt.Sprintf("Could not reach URL: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		sendError(w, resp.StatusCode, fmt.Sprintf("Target site returned an error: %s", resp.Status))
		return
	}

	// 4. Run Analysis
	// Pass the body to our parser (from parser.go)
	result, err := ParseHTML(resp.Body)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to parse HTML")
		return
	}

	// 5. Supplement with URL and Link checks (from checker.go)
	// For this task, we'd modify traverse slightly to return a list of links,
	// but for brevity, let's assume we extract them here or within ParseHTML.
	result.URL = req.URL

	// Example: Collecting all links and processing them
	result.Links = ProcessLinks(req.URL, result.discoveredLinks)

	// 6. Return JSON to React
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func sendError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(AnalysisResult{
		Error: &ErrorDetail{
			StatusCode: code,
			Message:    message,
		},
	})
}
