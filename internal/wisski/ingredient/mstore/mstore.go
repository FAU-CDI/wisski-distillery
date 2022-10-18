package mstore

import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/meta"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/tkw1536/goprogram/lib/collection"
)

// MStore implements metadata storage for this WissKI
type MStore struct {
	ingredient.Base
	*meta.Storage
}

// For is a Store for the provided value
type For[Value any] meta.Key

func (f For[Value]) Get(m *MStore) (value Value, err error) {
	err = m.Storage.Get(meta.Key(f), &value)
	return
}

func (f For[Value]) GetAll(m *MStore) (values []Value, err error) {
	err = m.Storage.GetAll(meta.Key(f), func(index, total int) any {
		if values == nil {
			values = make([]Value, total)
		}
		return &values[index]
	})
	return values, err
}

func (f For[Value]) Set(m *MStore, value Value) error {
	return m.Storage.Set(meta.Key(f), value)
}

func (f For[Value]) SetAll(m *MStore, values ...Value) error {
	return m.Storage.SetAll(meta.Key(f), collection.AsAny(values)...)
}

func (f For[Value]) Delete(m *MStore) error {
	return m.Storage.Delete(meta.Key(f))
}
