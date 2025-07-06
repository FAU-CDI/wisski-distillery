package cli

//spellchecker:words github goprogram meta
import (
	"go.tkw01536.de/goprogram"
	"go.tkw01536.de/goprogram/meta"
)

// Requirements are requirements for the WissKI Distillery.
type Requirements struct {
	// Do we need an installed distillery?
	NeedsDistillery bool

	// Automatically fail when cgo is enabled?
	FailOnCgo bool
}

// AllowsFlag checks if the provided flag may be passed to fullfill this requirement
// By default it is used only for help page generation, and may be inaccurate.
func (r Requirements) AllowsFlag(flag meta.Flag) bool {
	return r.NeedsDistillery
}

// Validate validates if this requirement is fullfilled for the provided global flags.
// It should return either nil, or an error of type exit.Error.
//
// Validate does not take into account AllowsOption, see ValidateAllowedOptions.
func (r Requirements) Validate(arguments goprogram.Arguments[Flags]) error {
	return nil
}
