package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SocketAdapter struct {
	Endpoint string
}

// func (s *SocketAdapter) Publish(result interface{}) {
// 	jsonData, _ := json.Marshal(result)
// 	resp, err := http.Post(s.Endpoint, "application/json", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		slog.Error("Failed to reach Socket Service: %v", "err", err)
// 		return
// 	}
// 	defer resp.Body.Close()
// }

func (s *SocketAdapter) Publish(result interface{}) error {
	// 1. Check for Marshalling errors first
	jsonData, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	// 2. Perform the POST request
	resp, err := http.Post(s.Endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		// This handles network-level errors (e.g., socket-service is down)
		return fmt.Errorf("failed to reach Socket Service: %w", err)
	}
	defer resp.Body.Close()

	// 3. Check for Non-200 Status Codes
	if resp.StatusCode >= 400 {
		return fmt.Errorf("socket service returned error status: %d", resp.StatusCode)
	}

	return nil
}
