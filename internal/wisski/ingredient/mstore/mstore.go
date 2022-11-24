package mstore

import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/meta"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
)

// MStore implements metadata storage for this WissKI
type MStore struct {
	ingredient.Base
	*meta.Storage
}

// For is a Store for the provided value
type For[Value any] meta.TypedKey[Value]

func (f For[Value]) Get(m *MStore) (value Value, err error) {
	return meta.TypedKey[Value](f).Get(m.Storage)
}

func (f For[Value]) GetAll(m *MStore) (values []Value, err error) {
	return meta.TypedKey[Value](f).GetAll(m.Storage)
}

func (f For[Value]) GetOrSet(m *MStore, dflt Value) (value Value, err error) {
	return meta.TypedKey[Value](f).GetOrSet(m.Storage, dflt)
}

func (f For[Value]) Set(m *MStore, value Value) error {
	return meta.TypedKey[Value](f).Set(m.Storage, value)
}

func (f For[Value]) SetAll(m *MStore, values ...Value) error {
	return meta.TypedKey[Value](f).SetAll(m.Storage, values...)
}

func (f For[Value]) Delete(m *MStore) error {
	return m.Storage.Delete(meta.Key(f))
}
