//spellchecker:words models
package models

var _ Model = Metadatum{}

// Metadatum represents a metadatum for a single model.
type Metadatum struct {
	Pk uint `gorm:"column:pk;primaryKey"`

	Key   string `gorm:"column:key;not null"` // key for the value, see the keys below
	Slug  string `gorm:"column:slug"`         // optional slug of instance
	Value []byte `gorm:"column:value"`        // serialized json value of the data
}

func (Metadatum) TableName() string {
	return "metadatum"
}
