// Package wdlog manages logging in the distillery
//
//spellchecker:words wdlog
package wdlog

//spellchecker:words context slog
import (
	"context"
	"io"
	"log/slog"
)

// New creates a new logger logging into the given output.
func New(out io.Writer, level slog.Level) *slog.Logger {
	return slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{
		Level: level,
	}))
}

// context.
type loggerKeyTyp struct{}

var loggerKey = loggerKeyTyp{}

// Set creates a new context that stores the given logger.
func Set(parent context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(parent, loggerKey, logger)
}

// If no logger is contained in the context, a no-op handler is returned.
func Of(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok && logger != nil {
		return logger
	}

	return slog.New(disabledHandler{})
}
