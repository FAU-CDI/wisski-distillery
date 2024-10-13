// Package wdlog manages logging in the distillery
package wdlog

import (
	"context"
	"io"
	"log/slog"
)

type Logger = *slog.Logger
type Level = slog.Level

// LevelFlag is a string representing a level.
// It can be used as a command-line flag.
type LevelFlag string

// Level returns the level belonging to this LevelFlag
func (ls LevelFlag) Level() Level {
	var level Level
	if err := level.UnmarshalText([]byte(ls)); err != nil {
		return slog.LevelInfo
	}
	return level
}

// New creates a new logger logging into the given output
func New(out io.Writer, level Level) Logger {
	return slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{
		Level: level,
	}))
}

// context
type loggerKeyTyp struct{}

var loggerKey = loggerKeyTyp{}

// Set creates a new context that stores the given logger
func Set(parent context.Context, logger Logger) context.Context {
	return context.WithValue(parent, loggerKey, logger)
}

// Of returns the logger stored in context.
// If no logger is contained in the context, a no-op handler is returned
func Of(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerKey).(Logger); ok && logger != nil {
		return logger
	}

	return slog.New(&disabledHandler{})
}

// disabledHandler implements slog.Handler as a no-op.
type disabledHandler struct{}

func (*disabledHandler) Enabled(context.Context, slog.Level) bool {
	return false
}

func (*disabledHandler) Handle(context.Context, slog.Record) error {
	panic("never called")
}

func (dh *disabledHandler) WithAttrs([]slog.Attr) slog.Handler {
	return dh
}

func (dh *disabledHandler) WithGroup(string) slog.Handler {
	return dh
}
