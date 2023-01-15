package barrel

import (
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/locker"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/mstore"
)

// Barrel provides access to the underlying Barrel
type Barrel struct {
	ingredient.Base
	Dependencies struct {
		Locker *locker.Locker
		MStore *mstore.MStore
	}
}

func (barrel *Barrel) DataPath() string {
	return filepath.Join(barrel.FilesystemBase, "data")
}
