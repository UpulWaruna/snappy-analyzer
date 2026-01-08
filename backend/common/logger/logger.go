package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
)

// Sensitive type masks values in logs automatically
type Sensitive string

func (s Sensitive) LogValue() slog.Value { return slog.StringValue("REDACTED") }

type Config struct {
	ServiceName string
	Level       slog.Level
	Writers     []io.Writer
}

// customHandler injects "source" only for Debug and Error levels
type customHandler struct{ slog.Handler }

func (h *customHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level == slog.LevelDebug || r.Level >= slog.LevelError {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		r.AddAttrs(slog.String("source", f.File))
	}
	return h.Handler.Handle(ctx, r)
}

func InitLogger(cfg Config) {
	var out io.Writer = os.Stdout
	if len(cfg.Writers) > 0 {
		out = io.MultiWriter(cfg.Writers...)
	}

	h := slog.NewJSONHandler(out, &slog.HandlerOptions{Level: cfg.Level})
	// Add the service name globally
	topHandler := h.WithAttrs([]slog.Attr{slog.String("service", cfg.ServiceName)})

	slog.SetDefault(slog.New(&customHandler{topHandler}))
}

// Scoped returns a logger pre-filled with the required Groups
func Scoped(userID, userName, reqID, url, method string) *slog.Logger {
	return slog.Default().With(
		slog.Group("userGroup", "id", userID, "name", userName),
		slog.Group("requestGroup", "id", reqID, "url", url, "type", method),
	)
}
