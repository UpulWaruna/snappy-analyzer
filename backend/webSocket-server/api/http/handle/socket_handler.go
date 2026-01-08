package handle

import (
	"common/logger"
	"encoding/json"
	"net/http"
	"webSocket-server/application/socket"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type SocketHandler struct {
	UseCase *socket.BroadcastUseCase
}

func (h *SocketHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
	l := logger.Scoped("socket", "user", "ws_conn", r.URL.Path, r.Method)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		l.Error("upgrade failed", "error", err)
		return
	}

	h.UseCase.Hub.Clients[ws] = true
	l.Info("Client connected")

	// Keep alive until client closes
	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			delete(h.UseCase.Hub.Clients, ws)
			l.Info("Client disconnected")
			break
		}
	}
}

func (h *SocketHandler) HandlePublish(w http.ResponseWriter, r *http.Request) {
	l := logger.Scoped("socket", "worker", "pub_req", r.URL.Path, r.Method)

	var data interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		l.Error("decode failed", "error", err)
		http.Error(w, "Invalid JSON", 400)
		return
	}

	h.UseCase.Publish(data)
	w.WriteHeader(http.StatusAccepted)
}
