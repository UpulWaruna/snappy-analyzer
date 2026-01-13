package socket

import (
	"log/slog"
	"webSocket-server/domain/model"
)

type BroadcastUseCase struct {
	Hub *model.Hub
}

// Run starts the infinite loop to listen for messages
func (uc *BroadcastUseCase) Run() {
	for {
		msg := <-uc.Hub.Broadcast
		for client := range uc.Hub.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				slog.Error("websocket write error", "error", err)
				client.Close()
				delete(uc.Hub.Clients, client)
			}
		}
	}
}

func (uc *BroadcastUseCase) Publish(data interface{}) {
	uc.Hub.Broadcast <- data
}
