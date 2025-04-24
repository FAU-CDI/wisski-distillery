package errwrap

import (
	"github.com/tkw1536/goprogram/exit"
)

// DeferWrap replaces *err with wrap.WrapError(*err) iff *err is not nil.
func DeferWrap(wrap exit.Error, err *error) {
	// TODO: Check if this is still needed
	if err == nil || *err == nil {
		return
	}
	*err = wrap.WrapError(*err)
}
