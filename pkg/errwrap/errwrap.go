package errwrap

import (
	"fmt"

	"github.com/tkw1536/goprogram/exit"
)

// DeferWrap replaces *err with an error that wraps both wrap and the underlying error.
// Deprecated: Manually wrap where needed.
func DeferWrap(wrap exit.Error, retval *error) {
	// TODO: Check if this is still needed
	if retval == nil || *retval == nil {
		return
	}
	*retval = fmt.Errorf("%w: %w", wrap, *retval)
}
