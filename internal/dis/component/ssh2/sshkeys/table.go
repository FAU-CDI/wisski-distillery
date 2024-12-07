//spellchecker:words sshkeys
package sshkeys

//spellchecker:words context reflect github wisski distillery internal component models gliderlabs
import (
	"context"
	"reflect"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/gliderlabs/ssh"
)

func (ssh2 *SSHKeys) TableInfo() component.TableInfo {
	return component.TableInfo{
		Model: reflect.TypeFor[models.Keys](),
		Name:  models.KeysTable,
	}
}

// Keys returns a list of keys for the given user
func (ssh2 *SSHKeys) Keys(ctx context.Context, user string) ([]models.Keys, error) {
	// the empty user has no key
	if user == "" {
		return nil, nil
	}

	// get the table
	table, err := ssh2.dependencies.SQL.QueryTable(ctx, ssh2)
	if err != nil {
		return nil, err
	}

	var keys []models.Keys

	// make a query to find all keys (in the underlying model)
	query := table.Find(&keys, &models.Keys{User: user})
	if query.Error != nil {
		return nil, query.Error
	}

	return keys, nil
}

// Add adds a new key to the given user, unless it already exists
func (ssh2 *SSHKeys) Add(ctx context.Context, user string, comment string, key ssh.PublicKey) error {
	// check that the given user exists
	{
		_, err := ssh2.dependencies.Auth.User(ctx, user)
		if err != nil {
			return err
		}
	}

	// fetch all the keys
	keys, err := ssh2.Keys(ctx, user)
	if err != nil {
		return err
	}

	pks := make([]ssh.PublicKey, 0, len(keys))
	for _, key := range keys {
		if pk := key.PublicKey(); pk != nil {
			pks = append(pks, pk)
		}
	}

	// key already exists
	if KeyOneOf(pks, key) {
		return nil
	}

	// create a new key with the given comment
	mk := models.Keys{
		User:    user,
		Comment: comment,
	}
	mk.SetPublicKey(key)

	// get the table
	table, err := ssh2.dependencies.SQL.QueryTable(ctx, ssh2)
	if err != nil {
		return err
	}

	// create the key instance
	return table.Create(&mk).Error
}

// Remove removes a given publuc key from a user.
func (ssh2 *SSHKeys) Remove(ctx context.Context, user string, key ssh.PublicKey) error {
	// find all the keys for the given user
	keys, err := ssh2.Keys(ctx, user)
	if err != nil {
		return err
	}

	// iterate and find all the public keys
	var pks []uint
	for _, candidate := range keys {
		if ssh.KeysEqual(candidate.PublicKey(), key) {
			pks = append(pks, candidate.Pk)
		}
	}

	// nothing to delete
	if len(pks) == 0 {
		return nil
	}

	// query the table again
	table, err := ssh2.dependencies.SQL.QueryTable(ctx, ssh2)
	if err != nil {
		return nil
	}

	// and do the delete
	return table.Where("pk in ?", pks).Delete(&models.Keys{}).Error
}

func (ssh2 *SSHKeys) OnUserDelete(ctx context.Context, user *models.User) error {
	// get the table
	table, err := ssh2.dependencies.SQL.QueryTable(ctx, ssh2)
	if err != nil {
		return err
	}

	// delete all keys for the user
	return table.Delete(&models.Keys{}, &models.Keys{User: user.User}).Error
}
