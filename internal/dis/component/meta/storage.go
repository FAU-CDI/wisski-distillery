package meta

import (
	"encoding/json"
	"errors"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"gorm.io/gorm"
)

// Key represents a key for metadata.
type Key string

// ErrMetadatumNotSet is returned by various [MetaStorage] functions when a metadatum is not set
var ErrMetadatumNotSet = errors.New("metadatum not set")

// Storage manages metadata for either the entire distillery, or a single slug
type Storage struct {
	Slug string
	sql  *sql.SQL
}

// Get retrieves metadata with the provided key and deserializes the first one into target.
// If no metadatum exists, returns [ErrMetadatumNotSet].
func (s Storage) Get(key Key, target any) error {
	table, err := s.sql.QueryTable(true, models.MetadataTable)
	if err != nil {
		return err
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
	return json.Unmarshal(datum.Value, target)
}

// GetAll receives all metadata with the provided keys.
// For each received value, the targets function is called with the current index, and total number of results.
// The function is intended to return a target for deserialization.
//
// When no metadatum exists, targets is not called, and nil error is returned.
func (s Storage) GetAll(key Key, target func(index, total int) any) error {
	table, err := s.sql.QueryTable(true, models.MetadataTable)
	if err != nil {
		return err
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
			return err
		}
	}

	return nil
}

// Delete deletes all metadata with the provided key.
func (s Storage) Delete(key Key) error {
	table, err := s.sql.QueryTable(true, models.MetadataTable)
	if err != nil {
		return err
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
func (s Storage) Set(key Key, value any) error {
	table, err := s.sql.QueryTable(true, models.MetadataTable)
	if err != nil {
		return err
	}

	// marshal the value
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return table.Transaction(func(tx *gorm.DB) error {
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
	})
}

// Set serializes values and stores them with the provided key.
// Any other metadata with the same key is deleted.
func (s Storage) SetAll(key Key, values ...any) error {
	table, err := s.sql.QueryTable(true, models.MetadataTable)
	if err != nil {
		return err
	}

	return table.Transaction(func(tx *gorm.DB) error {
		// delete the old values
		status := tx.Where(&models.Metadatum{Slug: s.Slug, Key: string(key)}).Delete(&models.Metadatum{})
		if err := status.Error; err != nil {
			return err
		}

		for _, value := range values {
			bytes, err := json.Marshal(value)
			if err != nil {
				return err
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
	})
}

// Purge removes all metadata, regardless of key.
func (s Storage) Purge() error {
	table, err := s.sql.QueryTable(true, models.MetadataTable)
	if err != nil {
		return err
	}

	status := table.Where("slug = ?", s.Slug).Delete(&models.Metadatum{})
	if status.Error != nil {
		return status.Error
	}
	return nil
}