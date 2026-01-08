package logger

import (
	"log/slog"
	"os"
)

// InitSharedLogger sets up a global JSON logger with a service name attribute
func InitSharedLogger(serviceName string) {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}).WithAttrs([]slog.Attr{
		slog.String("service", serviceName),
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
