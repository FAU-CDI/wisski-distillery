//spellchecker:words policy
package policy

//spellchecker:words context errors github wisski distillery internal models gorm clause
import (
	"context"
	"errors"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"gorm.io/gorm/clause"
)

var (
	ErrNoAccess = errors.New("no access")
	errInvalid  = errors.New("invalid parameters")
)

// Set sets a specific grant, overwriting any previous grant.
//
// User and Slug must not be empty.
// If DrupalUsername is empty, sets the username to be equal to the user.
func (policy *Policy) Set(ctx context.Context, grant models.Grant) error {
	if grant.DrupalUsername == "" {
		grant.DrupalUsername = grant.User
	}
	if grant.User == "" || grant.Slug == "" {
		return errInvalid
	}

	// check that the referenced user exists!
	{
		_, err := policy.dependencies.Auth.User(ctx, grant.User)
		if err != nil {
			return err
		}
	}

	// get the table
	table, err := policy.table(ctx)
	if err != nil {
		return err
	}

	// and create or update the given user / slug combination
	return table.Clauses(
		clause.OnConflict{OnConstraint: "user_slug", UpdateAll: true},
	).Create(&grant).Error
}

// Remove removes access for the given username form the given instance.
// The user not having access is not an error.
func (policy *Policy) Remove(ctx context.Context, username string, slug string) error {
	// empty username or slug never have acccess
	if username == "" || slug == "" {
		return errInvalid
	}

	// get the table
	table, err := policy.table(ctx)
	if err != nil {
		return err
	}

	// delete the access from the database
	return table.Delete(&models.Grant{}, models.Grant{User: username, Slug: slug}).Error
}

// User returns all grants for the given user.
func (policy *Policy) User(ctx context.Context, username string) (grants []models.Grant, err error) {
	if username == "" {
		return nil, errInvalid
	}

	// get the table
	table, err := policy.table(ctx)
	if err != nil {
		return nil, err
	}

	// find the grants
	err = table.Find(&grants, models.Grant{User: username}).Order("Slug asc").Error
	if err != nil {
		return nil, err
	}
	return grants, nil
}

// Instance returns all the grants for the given instance.
func (policy *Policy) Instance(ctx context.Context, slug string) (grants []models.Grant, err error) {
	if slug == "" {
		return nil, errInvalid
	}

	// get the table
	table, err := policy.table(ctx)
	if err != nil {
		return nil, err
	}

	// find the grants
	err = table.Find(&grants, models.Grant{Slug: slug}).Order("User asc").Error
	if err != nil {
		return nil, err
	}
	return grants, nil
}

// Has checks if the given username has access to the given instance.
// If the user has access, returns the provided grant.
//
// If the user does not have access, returns ErrNoAccess.
// Other errors may be returned in other cases.
func (policy *Policy) Has(ctx context.Context, username string, slug string) (grant models.Grant, err error) {
	// empty username or slug never have acccess
	if username == "" || slug == "" {
		return grant, errInvalid
	}

	// get the table
	table, err := policy.table(ctx)
	if err != nil {
		return grant, err
	}

	// read the access from the database
	res := table.Find(&grant, models.Grant{User: username, Slug: slug})
	if err := res.Error; err != nil {
		return grant, err
	}

	// if there were no rows affected, then there was no access granted
	if res.RowsAffected == 0 {
		return grant, ErrNoAccess
	}

	// return the username and admin
	return grant, nil
}
