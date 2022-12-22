package drush

import (
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/mstore"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
)

// Drush implements commands related to drush
type Drush struct {
	ingredient.Base
	Dependencies struct {
		Barrel *barrel.Barrel
		MStore *mstore.MStore
		PHP    *php.PHP
	}
}
