//spellchecker:words mstore
package mstore

//spellchecker:words context github wisski distillery internal component meta ingredient
import (
	"context"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/meta"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
)

// MStore implements metadata storage for this WissKI.
type MStore struct {
	ingredient.Base
	*meta.Storage
}

// For is a Store for the provided value.
type For[Value any] meta.TypedKey[Value]

func (f For[Value]) Get(ctx context.Context, m *MStore) (value Value, err error) {
	value, err = meta.TypedKey[Value](f).Get(ctx, m.Storage)
	if err != nil {
		return value, fmt.Errorf("failed to get values: %w", err)
	}
	return value, nil
}

func (f For[Value]) GetAll(ctx context.Context, m *MStore) (values []Value, err error) {
	values, err = meta.TypedKey[Value](f).GetAll(ctx, m.Storage)
	if err != nil {
		return nil, fmt.Errorf("failed to get all values: %w", err)
	}
	return values, nil
}

func (f For[Value]) GetOrSet(ctx context.Context, m *MStore, dflt Value) (value Value, err error) {
	value, err = meta.TypedKey[Value](f).GetOrSet(ctx, m.Storage, dflt)
	if err != nil {
		return value, fmt.Errorf("failed to get or set value: %w", err)
	}
	return value, nil
}

func (f For[Value]) Set(ctx context.Context, m *MStore, value Value) error {
	if err := meta.TypedKey[Value](f).Set(ctx, m.Storage, value); err != nil {
		return fmt.Errorf("failed to set value: %w", err)
	}
	return nil
}

func (f For[Value]) SetAll(ctx context.Context, m *MStore, values ...Value) error {
	if err := meta.TypedKey[Value](f).SetAll(ctx, m.Storage, values...); err != nil {
		return fmt.Errorf("failed to set values: %w", err)
	}
	return nil
}

func (f For[Value]) Delete(ctx context.Context, m *MStore) error {
	if err := m.Delete(ctx, meta.Key(f)); err != nil {
		return fmt.Errorf("failed to delete value: %w", err)
	}
	return nil
}
