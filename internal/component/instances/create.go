package instances

import (
	"errors"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/pkg/stringparser"
)

var errInvalidSlug = errors.New("not a valid slug")

// Create fills the struct for a new WissKI instance.
// It validates that slug is a valid name for an instance.
//
// It does not perform any checks if the instance already exists, or does the creation in the database.
func (instances *Instances) Create(slug string) (wissKI *wisski.WissKI, err error) {

	// make sure that the slug is valid!
	slug, err = stringparser.ParseSlug(instances.Environment, slug)
	if err != nil {
		return nil, errInvalidSlug
	}

	wissKI = new(wisski.WissKI)
	instances.use(wissKI)

	wissKI.Instance.Slug = slug
	wissKI.Instance.FilesystemBase = filepath.Join(instances.Path(), wissKI.Domain())

	wissKI.Instance.OwnerEmail = ""
	wissKI.Instance.AutoBlindUpdateEnabled = true

	// sql

	wissKI.Instance.SqlDatabase = instances.Config.MysqlDatabasePrefix + slug
	wissKI.Instance.SqlUsername = instances.Config.MysqlUserPrefix + slug

	wissKI.Instance.SqlPassword, err = instances.Config.NewPassword()
	if err != nil {
		return nil, err
	}

	// triplestore

	wissKI.Instance.GraphDBRepository = instances.Config.GraphDBRepoPrefix + slug
	wissKI.Instance.GraphDBUsername = instances.Config.GraphDBUserPrefix + slug

	wissKI.Instance.GraphDBPassword, err = instances.Config.NewPassword()
	if err != nil {
		return nil, err
	}

	// drupal

	wissKI.DrupalUsername = "admin" // TODO: Change this!

	wissKI.DrupalPassword, err = instances.Config.NewPassword()
	if err != nil {
		return nil, err
	}

	// store the instance in the object and return it!
	return wissKI, nil
}
