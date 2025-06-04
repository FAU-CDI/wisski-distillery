//spellchecker:words models
package models

var _ Model = Grant{}

// Grant represents an access grant to a specific user.
type Grant struct {
	Pk uint `gorm:"column:pk;primaryKey"`

	User string `gorm:"column:user;not null;index:user_slug,unique"`            // (distillery) username
	Slug string `gorm:"column:slug;not null;index:user_slug;index:drupal_slug"` // (distillery) instance slug

	DrupalUsername  string `gorm:"column:drupal_user;not null;index:drupal_slug,unique"` // drupal username
	DrupalAdminRole bool   `gorm:"column:admin;not null"`                                // drupal admin rights
}

func (Grant) TableName() string {
	return "grant"
}
