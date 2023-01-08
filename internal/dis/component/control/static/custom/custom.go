package custom

import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

// Custom implements theme and page customization.
type Custom struct {
	component.Base
	Dependencies struct {
		// nothing yet
	}
}

var (
	_ component.Backupable = (*Custom)(nil)
)
