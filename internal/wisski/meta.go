package wisski

import "github.com/FAU-CDI/wisski-distillery/internal/dis/component/meta"

func (wisski *WissKI) storage() *meta.Storage {
	return wisski.Meta.Storage(wisski.Slug)
}
