package main

import (
	"common/logger"
	"headlessBrowser-worker/adapter/external"
	"headlessBrowser-worker/api/http/handle"
	"headlessBrowser-worker/application/analysis"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	logger.InitLogger(logger.Config{ServiceName: "headless-worker", Level: slog.LevelDebug})

	// Dependency Manual Injection
	chrome := &external.ChromeAdapter{}
	publisher := &external.SocketAdapter{Endpoint: "http://socket-service:8081/publish"}
	useCase := &analysis.AnalyzeURLUseCase{Browser: chrome, Publisher: publisher}
	handler := &handle.AnalysisHandler{UseCase: useCase}

	mux := http.NewServeMux()
	mux.HandleFunc("/analyze", handler.HandleAnalyze)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	slog.Info("listening", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "err", err)
	}

}
