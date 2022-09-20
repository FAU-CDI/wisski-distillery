package instances

import (
	"encoding/json"
	"errors"

	"github.com/FAU-CDI/wisski-distillery/internal/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"gorm.io/gorm"
)

// MetaKey represents a key for metadata.
type MetaKey string

// ErrMetadatumNotSet is returned by various [MetaStorage] functions when a metadatum is not set
var ErrMetadatumNotSet = errors.New("metadatum not set")

// MetaStorage manages some metadata.
type MetaStorage interface {
	// Get retrieves metadata with the provided key and deserializes the first one into target.
	// If no metadatum exists, returns [ErrMetadatumNotSet].
	Get(key MetaKey, target any) error

	// GetAll receives all metadata with the provided keys.
	// For each received value, the targets function is called with the current index, and total number of results.
	// The function is intended to return a target for deserialization.
	//
	// When no metadatum exists, targets is not called, and nil error is returned.
	GetAll(key MetaKey, targets func(index, total int) any) error

	// Delete deletes all metadata with the provided key.
	Delete(key MetaKey) error

	// Set serializes value and stores it with the provided key.
	// Any other metadata with the same key is deleted.
	Set(key MetaKey, value any) error

	// Add serializes values and stores each as associated with the provided key.
	// Already existing metadata is left intact.
	Add(key MetaKey, values ...any) error

	// Purge removes all metadata, regardless of key.
	Purge() error
}

// Metadata returns a system-wide [MetaStorage].
func (instances *Instances) Metadata() MetaStorage {
	return &storage{
		SQL:  instances.SQL,
		Slug: "", // not associated to any slug
	}
}

// Metadata returns a [MetaStorage] that manages metadata related to this WissKI instance.
// It will be automatically deleted once the instance is deleted.
func (wisski *WissKI) Metadata() MetaStorage {
	return &storage{
		SQL:  wisski.instances.SQL,
		Slug: wisski.Slug, // associated to this instance
	}
}

// storage implements MetaStorage
type storage struct {
	SQL  *sql.SQL
	Slug string
}

func (s *storage) Get(key MetaKey, target any) error {
	table, err := s.SQL.QueryTable(true, models.MetadataTable)
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

func (s *storage) GetAll(key MetaKey, target func(index, total int) any) error {
	table, err := s.SQL.QueryTable(true, models.MetadataTable)
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

func (s *storage) Delete(key MetaKey) error {
	table, err := s.SQL.QueryTable(true, models.MetadataTable)
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

func (s *storage) Set(key MetaKey, value any) error {
	table, err := s.SQL.QueryTable(true, models.MetadataTable)
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

func (s *storage) Add(key MetaKey, values ...any) error {
	table, err := s.SQL.QueryTable(true, models.MetadataTable)
	if err != nil {
		return err
	}

	return table.Transaction(func(tx *gorm.DB) error {
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

func (s *storage) Purge() error {
	table, err := s.SQL.QueryTable(true, models.MetadataTable)
	if err != nil {
		return err
	}

	status := table.Where("slug = ?", s.Slug).Delete(&models.Metadatum{})
	if status.Error != nil {
		return status.Error
	}
	return nil
}
