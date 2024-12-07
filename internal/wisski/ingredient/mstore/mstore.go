//spellchecker:words mstore
package mstore

//spellchecker:words context github wisski distillery internal component meta ingredient
import (
	"context"

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

func (f For[Value]) Get(ctx context.Context, m *MStore) (value Value, err error) {
	return meta.TypedKey[Value](f).Get(ctx, m.Storage)
}

func (f For[Value]) GetAll(ctx context.Context, m *MStore) (values []Value, err error) {
	return meta.TypedKey[Value](f).GetAll(ctx, m.Storage)
}

func (f For[Value]) GetOrSet(ctx context.Context, m *MStore, dflt Value) (value Value, err error) {
	return meta.TypedKey[Value](f).GetOrSet(ctx, m.Storage, dflt)
}

func (f For[Value]) Set(ctx context.Context, m *MStore, value Value) error {
	return meta.TypedKey[Value](f).Set(ctx, m.Storage, value)
}

func (f For[Value]) SetAll(ctx context.Context, m *MStore, values ...Value) error {
	return meta.TypedKey[Value](f).SetAll(ctx, m.Storage, values...)
}

func (f For[Value]) Delete(ctx context.Context, m *MStore) error {
	return m.Storage.Delete(ctx, meta.Key(f))
}
