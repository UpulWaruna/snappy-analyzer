package main

import (
	"common/logger"
	"headlessBrowser-worker/adapter/external"
	"headlessBrowser-worker/adapter/repository/memory"
	"headlessBrowser-worker/api/http/handle"
	"headlessBrowser-worker/api/middleware"
	"headlessBrowser-worker/application/analysis"
	"headlessBrowser-worker/domain/service"
	"log/slog"
	"net/http"
	"time"

	"github.com/rs/cors"
)

func main() {
	logger.InitLogger(logger.Config{ServiceName: "headless-worker", Level: slog.LevelDebug})

	// Dependency Manual Injection
	chrome := &external.ChromeAdapter{}
	publisher := &external.SocketAdapter{Endpoint: "http://socket-service:8081/publish"}
	repo := memory.NewInMemoryRepository()
	dataService := service.NewAnalysisDataService(repo)
	useCase := &analysis.AnalyzeURLUseCase{Browser: chrome, Publisher: publisher, DataService: dataService}
	handler := &handle.AnalysisHandler{UseCase: useCase}

	mux := http.NewServeMux()
	mux.HandleFunc("/analyze", handler.HandleAnalyze)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })

	// Setup CORS options
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Your React URL
		AllowedMethods: []string{"POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Idempotency-Key"},
		Debug:          true, // Useful for troubleshooting
	})

	// Wrap middleware
	handlerChain :=
		c.Handler(
			middleware.LoggingMiddleware(
				middleware.IdempotencyMiddleware(
					middleware.CheckInCacheMiddleware(dataService)(mux),
				),
			),
		)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handlerChain,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	//srv.Handler = c.Handler(finalHandler)

	slog.Info("listening", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "err", err)
	}

}
