package custom

import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
)

// Custom implements theme and page customization.
type Custom struct {
	component.Base
	Dependencies struct {
		Routeables []component.Routeable
		Menuable   []component.Menuable
	}
	menu lazy.Lazy[[]component.MenuItem]
}

var (
	_ component.Backupable = (*Custom)(nil)
	_ component.Menuable   = (*Custom)(nil)
)
