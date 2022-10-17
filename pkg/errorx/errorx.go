package errorx

import "github.com/tkw1536/goprogram/lib/collection"

// First returns the first non-nil error, or nil otherwise.
func First(errors ...error) error {
	return collection.First(errors, func(err error) bool { return err != nil })
}
