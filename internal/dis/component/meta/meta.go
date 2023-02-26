package meta

import (
	"sync"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/tkw1536/pkglib/reflectx"
)

// Component meta is responsible for managing metadata per WissKI Instance
type Meta struct {
	component.Base
	Dependencies struct {
		SQL *sql.SQL
	}

	sl sync.Mutex
	sc map[string]*Storage
}

var (
	_ component.Provisionable = (*Meta)(nil)
	_ component.Table         = (*Meta)(nil)
)

func (*Meta) TableInfo() component.TableInfo {
	return component.TableInfo{
		Model: reflectx.TypeOf[models.Metadatum](),
		Name:  models.MetadataTable,
	}
}

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
		Slug:  slug,
		sql:   meta.Dependencies.SQL,
		table: meta,
	}
	return meta.sc[slug]
}
