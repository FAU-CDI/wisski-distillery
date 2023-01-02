package models

// GrantTable is the name of the table the 'Grant' model is stored in.
const GrantTable = "grant"

// Grant represents an access grant to a specific user
type Grant struct {
	Pk uint `gorm:"column:pk;primaryKey"`

	User string `gorm:"column:user;not null;uniqueIndex:user_slug"` // (distillery) username
	Slug string `gorm:"column:slug;not null;uniqueIndex:user_slug"` // (distillery) instance slug

	DrupalUsername  string `gorm:"column:drupal_user;not null"` // drupal username
	DrupalAdminRole bool   `gorm:"column:admin;not null"`       // drupal admin rights
}
