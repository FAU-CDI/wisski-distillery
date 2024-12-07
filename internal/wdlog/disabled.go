//spellchecker:words wdlog
package wdlog

//spellchecker:words context slog
import (
	"context"
	"log/slog"
)

// disabledHandler implements slog.Handler as a no-op.
type disabledHandler struct{}

var _ slog.Handler = (*disabledHandler)(nil)

func (disabledHandler) Enabled(context.Context, slog.Level) bool {
	return false
}

func (disabledHandler) Handle(context.Context, slog.Record) error {
	panic("never called")
}

func (dh disabledHandler) WithAttrs([]slog.Attr) slog.Handler {
	return dh
}

func (dh disabledHandler) WithGroup(string) slog.Handler {
	return dh
}
