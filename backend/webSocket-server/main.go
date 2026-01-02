package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"common/logger"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow React frontend
}

// Hub maintains the set of active clients
type Hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan interface{}
}

var hub = Hub{
	clients:   make(map[*websocket.Conn]bool),
	broadcast: make(chan interface{}),
}

func main() {
	// Start the broadcaster in a background goroutine
	go handleMessages()

	// Route for React to connect
	http.HandleFunc("/ws", handleConnections)

	// Route for Backend Worker to "Push" data
	http.HandleFunc("/publish", handlePublish)

	// ADD THIS LINE:
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	logger.InitSharedLogger("webSocket Service")
	slog.Info("Socket Service started on", "port", 8081)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	hub.clients[ws] = true
	log.Println("New Client Connected")

	// Keep connection alive
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			delete(hub.clients, ws)
			break
		}
	}
}

func handlePublish(w http.ResponseWriter, r *http.Request) {
	var analysisResult interface{}
	if err := json.NewDecoder(r.Body).Decode(&analysisResult); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Send to the broadcast channel
	hub.broadcast <- analysisResult
	w.WriteHeader(http.StatusAccepted)
}

func handleMessages() {
	for {
		msg := <-hub.broadcast
		for client := range hub.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(hub.clients, client)
			}
		}
	}
}
