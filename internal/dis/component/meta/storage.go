//spellchecker:words meta
package meta

//spellchecker:words context encoding json errors github wisski distillery internal component models pkglib collection gorm
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/tkw1536/pkglib/collection"
	"gorm.io/gorm"
)

// Key represents a key for metadata.
type Key string

// ErrMetadatumNotSet is returned by various [MetaStorage] functions when a metadatum is not set.
var ErrMetadatumNotSet = errors.New("metadatum not set")

// Storage manages metadata for either the entire distillery, or a single slug.
type Storage struct {
	Slug string

	table component.Table
	sql   *sql.SQL
}

// Get retrieves metadata with the provided key and deserializes the first one into target.
// If no metadatum exists, returns [ErrMetadatumNotSet].
func (s Storage) Get(ctx context.Context, key Key, target any) error {
	table, err := s.sql.QueryTable(ctx, s.table)
	if err != nil {
		return fmt.Errorf("failed to query table: %w", err)
	}

	// read the datum from the database
	var datum models.Metadatum
	status := table.Where(&models.Metadatum{Slug: s.Slug, Key: string(key)}).Order("pk DESC").Find(&datum)

	// check if there was an error
	if err := status.Error; err != nil {
		return err
	}

	// check if e actually found it!
	if status.RowsAffected == 0 {
		return ErrMetadatumNotSet
	}

	// and do the unmarshaling!
	if err := json.Unmarshal(datum.Value, target); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}
	return nil
}

// GetAll receives all metadata with the provided keys.
// For each received value, the targets function is called with the current index, and total number of results.
// The function is intended to return a target for deserialization.
//
// When no metadatum exists, targets is not called, and nil error is returned.
func (s Storage) GetAll(ctx context.Context, key Key, target func(index, total int) any) error {
	table, err := s.sql.QueryTable(ctx, s.table)
	if err != nil {
		return fmt.Errorf("failed to query table: %w", err)
	}

	// read the datum from the database
	var data []models.Metadatum
	status := table.Where(&models.Metadatum{Slug: s.Slug, Key: string(key)}).Find(&data)

	// check if there was an error
	if err := status.Error; err != nil {
		return err
	}

	// unpack all of them into the destination
	for index, datum := range data {
		err := json.Unmarshal(datum.Value, target(index, len(data)))
		if err != nil {
			return fmt.Errorf("failed to unmarshal: %w", err)
		}
	}

	return nil
}

// Delete deletes all metadata with the provided key.
func (s Storage) Delete(ctx context.Context, key Key) error {
	table, err := s.sql.QueryTable(ctx, s.table)
	if err != nil {
		return fmt.Errorf("failed to query table: %w", err)
	}

	// delete all the values
	status := table.Where(&models.Metadatum{Slug: s.Slug, Key: string(key)}).Delete(&models.Metadatum{})
	if err := status.Error; err != nil {
		return err
	}
	return nil
}

// Set serializes value and stores it with the provided key.
// Any other metadata with the same key is deleted.
func (s Storage) Set(ctx context.Context, key Key, value any) error {
	// marshal the value
	bytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	table, err := s.sql.QueryTable(ctx, s.table)
	if err != nil {
		return fmt.Errorf("failed to query table: %w", err)
	}

	if err := table.Transaction(func(tx *gorm.DB) error {
		// delete the old values
		status := tx.Where(&models.Metadatum{Slug: s.Slug, Key: string(key)}).Delete(&models.Metadatum{})
		if err := status.Error; err != nil {
			return err
		}

		// create the new item to insert
		status = tx.Create(&models.Metadatum{
			Key:   string(key),
			Slug:  s.Slug,
			Value: bytes,
		})
		if err := status.Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	return nil
}

// Set serializes values and stores them with the provided key.
// Any other metadata with the same key is deleted.
func (s Storage) SetAll(ctx context.Context, key Key, values ...any) error {
	table, err := s.sql.QueryTable(ctx, s.table)
	if err != nil {
		return fmt.Errorf("failed to query table: %w", err)
	}

	if err := table.Transaction(func(tx *gorm.DB) error {
		// delete the old values
		status := tx.Where(&models.Metadatum{Slug: s.Slug, Key: string(key)}).Delete(&models.Metadatum{})
		if err := status.Error; err != nil {
			return err
		}

		for _, value := range values {
			bytes, err := json.Marshal(value)
			if err != nil {
				return fmt.Errorf("failed to marshal value: %w", err)
			}

			// create the new item to insert
			status := tx.Create(&models.Metadatum{
				Key:   string(key),
				Slug:  s.Slug,
				Value: bytes,
			})
			if err := status.Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("transation failed: %w", err)
	}
	return nil
}

// Purge removes all metadata, regardless of key.
func (s Storage) Purge(ctx context.Context) error {
	table, err := s.sql.QueryTable(ctx, s.table)
	if err != nil {
		return fmt.Errorf("failed to query table: %w", err)
	}

	status := table.Where("slug = ?", s.Slug).Delete(&models.Metadatum{})
	if status.Error != nil {
		return status.Error
	}
	return nil
}

// TypedKey represents a convenience wrapper for a given with a given value.
type TypedKey[Value any] Key

func (f TypedKey[Value]) Get(ctx context.Context, s *Storage) (value Value, err error) {
	err = s.Get(ctx, Key(f), &value)
	return
}

func (f TypedKey[Value]) GetOrSet(ctx context.Context, s *Storage, dflt Value) (value Value, err error) {
	value, err = f.Get(ctx, s)
	if errors.Is(err, ErrMetadatumNotSet) {
		value = dflt
		err = f.Set(ctx, s, value)
	}
	return
}

func (f TypedKey[Value]) GetAll(ctx context.Context, m *Storage) (values []Value, err error) {
	err = m.GetAll(ctx, Key(f), func(index, total int) any {
		if values == nil {
			values = make([]Value, total)
		}
		return &values[index]
	})
	return values, err
}

func (f TypedKey[Value]) Set(ctx context.Context, m *Storage, value Value) error {
	return m.Set(ctx, Key(f), value)
}

func (f TypedKey[Value]) SetAll(ctx context.Context, m *Storage, values ...Value) error {
	return m.SetAll(ctx, Key(f), collection.AsAny(values)...)
}

func (f TypedKey[Value]) Delete(ctx context.Context, m *Storage) error {
	return m.Delete(ctx, Key(f))
}
