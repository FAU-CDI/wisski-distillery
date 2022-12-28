package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// ErrUserNotFound is returned when a user is not found
var ErrUserNotFound = errors.New("user not found")

// Users returns all users in the database
func (auth *Auth) Users(ctx context.Context) (users []*AuthUser, err error) {
	// query the user table
	table, err := auth.Dependencies.SQL.QueryTable(ctx, false, models.UserTable)
	if err != nil {
		return
	}

	// find all the users
	var dUsers []models.User
	err = table.Find(&dUsers).Error
	if err != nil {
		return nil, err
	}

	// and map them to high-level user objects
	users = make([]*AuthUser, len(dUsers))
	for i, user := range dUsers {
		users[i] = &AuthUser{
			User: user,
			auth: auth,
		}
	}

	return users, nil
}

// User returns a single user.
// If the user does not exist, returns ErrUserNotFound.
func (auth *Auth) User(ctx context.Context, name string) (user *AuthUser, err error) {
	// quick and dirty check for the empty username (which is not allowed)
	if name == "" {
		return nil, ErrUserNotFound
	}

	// return the user
	table, err := auth.Dependencies.SQL.QueryTable(ctx, false, models.UserTable)
	if err != nil {
		return
	}

	user = &AuthUser{}

	// find the user
	res := table.Where(&models.User{User: name}).Find(&user.User)
	err = res.Error
	if err != nil {
		return
	}

	// check if the user was not found
	if res.RowsAffected == 0 {
		return nil, ErrUserNotFound
	}

	user.auth = auth

	return
}

// CreateUser creates a new user and returns it.
// The user is not associated to any WissKIs, and has no password set.
func (auth *Auth) CreateUser(ctx context.Context, name string) (user *AuthUser, err error) {
	// return the user
	table, err := auth.Dependencies.SQL.QueryTable(ctx, false, models.UserTable)
	if err != nil {
		return
	}

	user = &AuthUser{
		User: models.User{
			User:    name,
			Enabled: false,
		},
	}

	// do the create statement
	err = table.Create(&user.User).Error
	if err != nil {
		return nil, err
	}

	user.auth = auth
	return user, nil
}

// AuthUser represents an authorized user
type AuthUser struct {
	auth *Auth
	models.User
}

func (au *AuthUser) String() string {
	if au == nil {
		return "User{nil}"
	}
	hasPassword := len(au.PasswordHash) > 0
	return fmt.Sprintf("User{Name:%q,Enabled:%t,HasPassword:%t,Admin:%t}", au.User.User, au.User.Enabled, hasPassword, au.User.Admin)
}

// SetPassword sets the password for this user and turns the user on
func (au *AuthUser) SetPassword(ctx context.Context, password []byte) (err error) {
	au.User.PasswordHash, err = bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	au.User.Enabled = true
	return au.Save(ctx)
}

// UnsetPassword removes the password from this user, and disables them
func (au *AuthUser) UnsetPassword(ctx context.Context) error {
	au.User.PasswordHash = nil
	au.User.Enabled = false
	return au.Save(ctx)
}

var ErrNoUser = errors.New("user is nil")
var ErrUserDisabled = errors.New("user is disabled")
var ErrUserBlank = errors.New("user has no password set")

// CheckPassword checks if this user can login with the provided password.
// Returns nil on success, an error otherwise.
func (au *AuthUser) CheckPassword(ctx context.Context, password []byte) error {
	if au == nil {
		return ErrNoUser
	}
	if !au.User.Enabled {
		return ErrUserDisabled
	}

	if len(au.User.PasswordHash) == 0 {
		return ErrUserDisabled
	}

	return bcrypt.CompareHashAndPassword(au.User.PasswordHash, password)
}

// MakeAdmin makes this user an admin, and saves the update in the database.
// If the user is already an admin, does not return an error.
func (au *AuthUser) MakeAdmin(ctx context.Context) error {
	au.User.Admin = true
	return au.Save(ctx)
}

// MakeRegular removes admin rights from this user.
// If this user is not an dmin, does not return an error.
func (au *AuthUser) MakeRegular(ctx context.Context) error {
	au.User.Admin = true
	return au.Save(ctx)
}

// Save saves the given user in the database
func (au *AuthUser) Save(ctx context.Context) error {
	table, err := au.auth.Dependencies.SQL.QueryTable(ctx, false, models.UserTable)
	if err != nil {
		return err
	}
	return table.Save(&au.User).Error
}

// Delete deletes the user from the database
func (au *AuthUser) Delete(ctx context.Context) error {
	table, err := au.auth.Dependencies.SQL.QueryTable(ctx, false, models.UserTable)
	if err != nil {
		return err
	}
	return table.Delete(&au.User).Error
}
