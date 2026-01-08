package main

import (
	"common/logger"
	"log/slog"
	"net/http"
	"webSocket-server/api/http/handle"
	"webSocket-server/application/socket"
	"webSocket-server/domain/model"
)

func main() {
	// 1. Setup Logger
	logger.InitLogger(logger.Config{
		ServiceName: "socket-server",
		Level:       slog.LevelDebug,
	})

	// 2. Initialize DDD Layers
	hub := model.NewHub()
	useCase := &socket.BroadcastUseCase{Hub: hub}
	handler := &handle.SocketHandler{UseCase: useCase}

	// 3. Start Background Broadcaster
	go useCase.Run()

	// 4. Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handler.HandleWS)
	mux.HandleFunc("/publish", handler.HandlePublish)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	slog.Info("Socket server starting", "port", 8081)
	if err := http.ListenAndServe(":8081", mux); err != nil {
		slog.Error("server failed", "error", err)
	}
}
