package external

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
)

type SocketAdapter struct {
	Endpoint string
}

func (s *SocketAdapter) Publish(result interface{}) {
	jsonData, _ := json.Marshal(result)
	resp, err := http.Post(s.Endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Info("Failed to reach Socket Service: %v", err)
		return
	}
	defer resp.Body.Close()
}
