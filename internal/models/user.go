package models

// UserTable is the name of the table the [`User`] model is stored in.
const UserTable = "users"

// User represents a distillery user
type User struct {
	Pk uint `gorm:"column:pk;primaryKey"`

	User         string `gorm:"column:user;not null;unique"` // name of the user
	PasswordHash []byte `gorm:"column:password"`             // password of the user, hashed

	TOTPEnabled *bool  `gorm:"column:totpenabled"` // is totp enabled for the user
	TOTPURL     string `gorm:"column:totp"`        // the totp of the user

	Enabled *bool `gorm:"enabled;not null"`
	Admin   *bool `gorm:"column:admin;not null"`
}

func (user *User) IsAdmin() bool {
	return user.Admin != nil && *user.Admin
}

func (user *User) SetAdmin(v bool) {
	user.Admin = &v
}

func (user *User) IsEnabled() bool {
	return user.Enabled != nil && *user.Enabled
}

func (user *User) SetEnabled(v bool) {
	user.Enabled = &v
}

func (user *User) IsTOTPEnabled() bool {
	return user.TOTPEnabled != nil && *user.TOTPEnabled
}

func (user *User) SetTOTPEnabled(v bool) {
	user.TOTPEnabled = &v
}
