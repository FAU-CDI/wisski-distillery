package models

// UserTable is the name of the table the [`User`] model is stored in.
const UserTable = "users"

// User represents a distillery user
type User struct {
	Pk uint `gorm:"column:pk;primaryKey"`

	User         string `gorm:"column:user;not null;unique"` // name of the user
	PasswordHash []byte `gorm:"column:password"`             // password of the user, hashed

	Enabled bool `gorm:"enabled;not null"`
	Admin   bool `gorm:"column:admin;not null"`
}
