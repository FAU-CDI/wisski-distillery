package barrel

import (
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/locker"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/mstore"
)

// Barrel provides access to the underlying Barrel
type Barrel struct {
	ingredient.Base
	dependencies struct {
		Locker *locker.Locker
		MStore *mstore.MStore
	}
}

const (
	BaseDirectory     = "/var/www/data"
	ComposerDirectory = BaseDirectory + "/project"
	WebDirectory      = ComposerDirectory + "/web"
	OntologyDirectory = SitesDirectory + "/default/files/ontology"
	SitesDirectory    = WebDirectory + "/sites"
	WissKIDirectory   = WebDirectory + "/modules/contrib/wisski"
)
