package cli

// Requirements are requirements for the WissKI Distillery.
type Requirements struct {
	// Do we need an installed distillery?
	NeedsDistillery bool

	// Automatically fail when cgo is enabled?
	FailOnCgo bool
}
