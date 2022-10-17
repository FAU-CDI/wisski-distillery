package meta

import (
	"sync"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/sql"
)

// Component meta is responsible for managing metadata per WissKI Instance
type Meta struct {
	component.ComponentBase

	SQL *sql.SQL

	sl sync.Mutex
	sc map[string]*Storage
}

func (*Meta) Name() string { return "metadata" }

// Storage returns a Storage for the instance with the given slug.
// When slug is nil, returns a global storage.
func (meta *Meta) Storage(slug string) *Storage {
	meta.sl.Lock()
	defer meta.sl.Unlock()

	// create the cache (unless it already exists)
	if meta.sc == nil {
		meta.sc = make(map[string]*Storage)
	}

	// cache hit
	if storage, ok := meta.sc[slug]; ok {
		return storage
	}

	// create a new storage
	meta.sc[slug] = &Storage{
		Slug: slug,
		sql:  meta.SQL,
	}
	return meta.sc[slug]
}
