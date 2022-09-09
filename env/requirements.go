package env

import (
	"github.com/tkw1536/goprogram"
	"github.com/tkw1536/goprogram/meta"
)

type Requirements struct {
	// Do we need an installed distillery?
	NeedsDistillery bool
}

// AllowsFlag checks if the provided flag may be passed to fullfill this requirement
// By default it is used only for help page generation, and may be inaccurate.
func (r Requirements) AllowsFlag(flag meta.Flag) bool {
	return true
}

// Validate validates if this requirement is fullfilled for the provided global flags.
// It should return either nil, or an error of type exit.Error.
//
// Validate does not take into account AllowsOption, see ValidateAllowedOptions.
func (r Requirements) Validate(arguments goprogram.Arguments[struct{}]) error {
	return nil
}
