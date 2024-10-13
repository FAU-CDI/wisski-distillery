// Package wdlog manages logging in the distillery
package wdlog

import (
	"context"
	"io"

	"github.com/rs/zerolog"
)

type Logger = zerolog.Logger
type Level = zerolog.Level

// LevelFlag is a string representing a level.
// It can be used as a command-line flag.
type LevelFlag string

// Level returns the level belonging to this LevelFlag
func (ls LevelFlag) Level() Level {
	level, err := zerolog.ParseLevel(string(ls))
	if err != nil {
		return zerolog.InfoLevel
	}
	return level
}

// New creates a new logger logging into the given output
func New(out io.Writer, level Level) Logger {
	writer := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.Out = out
	})

	return zerolog.New(writer).Level(level)
}

// Set creates a new context that stores the given logger
func Set(parent context.Context, logger *Logger) context.Context {
	return logger.WithContext(parent)
}

// Of returns the logger stored in context.
// If no logger is in the given context, returns a no-op logger.
func Of(ctx context.Context) *Logger {
	return zerolog.Ctx(ctx)
}
