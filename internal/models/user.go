//spellchecker:words models
package models

var _ Model = User{}

// User represents a distillery user.
//
//nolint:recvcheck
type User struct {
	Pk uint `gorm:"column:pk;primaryKey"`

	User string `gorm:"column:user;not null;unique"` // name of the user

	PasswordHash []byte `gorm:"column:password"    json:"-"` // password of the user, hashed
	TOTPEnabled  *bool  `gorm:"column:totpenabled" json:"-"` // is totp enabled for the user
	TOTPURL      string `gorm:"column:totp"        json:"-"` // the totp of the user

	Enabled *bool `gorm:"enabled;not null"      json:"enabled"`
	Admin   *bool `gorm:"column:admin;not null" json:"admin"`
}

func (User) TableName() string {
	return "users"
}

func (user *User) HasPassword() bool {
	return len(user.PasswordHash) != 0
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
