package instances

import (
	"errors"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/pkg/stringparser"
)

var errInvalidSlug = errors.New("not a valid slug")

// Create fills the struct for a new WissKI instance.
// It validates that slug is a valid name for an instance.
//
// It does not perform any checks if the instance already exists, or does the creation in the database.
func (instances *Instances) Create(slug string) (wisski WissKI, err error) {

	// make sure that the slug is valid!
	if _, err := stringparser.ParseSlug(slug); err != nil {
		return wisski, errInvalidSlug
	}

	wisski.Instance.Slug = slug
	wisski.Instance.FilesystemBase = filepath.Join(instances.Dir, slug)

	wisski.Instance.OwnerEmail = ""
	wisski.Instance.AutoBlindUpdateEnabled = true

	// sql

	wisski.Instance.SqlDatabase = instances.Config.MysqlDatabasePrefix + slug
	wisski.Instance.SqlUsername = instances.Config.MysqlUserPrefix + slug

	wisski.Instance.SqlPassword, err = instances.Config.NewPassword()
	if err != nil {
		return WissKI{}, err
	}

	// triplestore

	wisski.Instance.GraphDBRepository = instances.Config.GraphDBRepoPrefix + slug
	wisski.Instance.GraphDBUsername = instances.Config.GraphDBUserPrefix + slug

	wisski.Instance.GraphDBPassword, err = instances.Config.NewPassword()
	if err != nil {
		return WissKI{}, err
	}

	// drupal

	wisski.DrupalUsername = "admin" // TODO: Change this!

	wisski.DrupalPassword, err = instances.Config.NewPassword()
	if err != nil {
		return wisski, err
	}

	// store the instance in the object and return it!
	wisski.instances = instances
	return wisski, nil
}
