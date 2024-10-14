package wdlog

import "log/slog"

// Flag is used to represents a level command line flag.
type Flag string

// Level parses the string into a log level.
// If it cannot be parsed, uses the default info level.
func (ls Flag) Level() (level slog.Level) {
	if err := level.UnmarshalText([]byte(ls)); err != nil {
		var zero slog.Level
		return zero
	}
	return
}
