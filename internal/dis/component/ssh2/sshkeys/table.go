//spellchecker:words sshkeys
package sshkeys

//spellchecker:words context reflect github wisski distillery internal component models gliderlabs
import (
	"context"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/gliderlabs/ssh"
)

func (ssh2 *SSHKeys) TableInfo() component.TableInfo {
	return component.TableInfo{
		Model: models.Keys{},
	}
}

// Keys returns a list of keys for the given user.
func (ssh2 *SSHKeys) Keys(ctx context.Context, user string) ([]models.Keys, error) {
	// the empty user has no key
	if user == "" {
		return nil, nil
	}

	// get the table
	table, err := ssh2.dependencies.SQL.OpenTable(ctx, ssh2)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}

	var keys []models.Keys

	// make a query to find all keys (in the underlying model)
	query := table.Find(&keys, &models.Keys{User: user})
	if query.Error != nil {
		return nil, fmt.Errorf("failed to find user: %w", query.Error)
	}

	return keys, nil
}

// Add adds a new key to the given user, unless it already exists.
func (ssh2 *SSHKeys) Add(ctx context.Context, user string, comment string, key ssh.PublicKey) error {
	// check that the given user exists
	{
		_, err := ssh2.dependencies.Auth.User(ctx, user)
		if err != nil {
			return fmt.Errorf("failed to authenticate user: %w", err)
		}
	}

	// fetch all the keys
	keys, err := ssh2.Keys(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to retrieve keys: %w", err)
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
	table, err := ssh2.dependencies.SQL.OpenTable(ctx, ssh2)
	if err != nil {
		return fmt.Errorf("failed to query ssh key table: %w", err)
	}

	// create the key instance
	if err := table.Create(&mk).Error; err != nil {
		return fmt.Errorf("failed to insert key into table: %w", err)
	}
	return nil
}

// Remove removes a given publuc key from a user.
func (ssh2 *SSHKeys) Remove(ctx context.Context, user string, key ssh.PublicKey) error {
	// find all the keys for the given user
	keys, err := ssh2.Keys(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to get keys: %w", err)
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
	table, err := ssh2.dependencies.SQL.OpenTable(ctx, ssh2)
	if err != nil {
		return nil
	}

	// and do the delete
	if err := table.Where("pk in ?", pks).Delete(&models.Keys{}).Error; err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}
	return nil
}

func (ssh2 *SSHKeys) OnUserDelete(ctx context.Context, user *models.User) error {
	// get the table
	table, err := ssh2.dependencies.SQL.OpenTable(ctx, ssh2)
	if err != nil {
		return fmt.Errorf("failkd to query user table: %w", err)
	}

	// delete all keys for the user
	if err := table.Delete(&models.Keys{}, &models.Keys{User: user.User}).Error; err != nil {
		return fmt.Errorf("faile to delete ssh keys for user: %w", err)
	}
	return nil
}
