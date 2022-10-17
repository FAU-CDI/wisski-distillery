package wisski

import "github.com/FAU-CDI/wisski-distillery/internal/component/meta"

func (wisski *WissKI) storage() *meta.Storage {
	return wisski.Meta.Storage(wisski.Slug)
}
