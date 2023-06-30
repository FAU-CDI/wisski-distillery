package instances

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/config/validators"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
)

var (
	errInvalidSlug    = errors.New("not a valid slug")
	errRestrictedSlug = errors.New("restricted slug")
)

// Create fills the struct for a new WissKI instance.
// It validates that slug is a valid name for an instance.
//
// It does not perform any checks if the instance already exists, or does the creation in the database.
func (instances *Instances) Create(slug string, phpversion string) (wissKI *wisski.WissKI, err error) {

	// make sure that the slug is valid!
	slug, err = instances.IsValidSlug(slug)
	if err != nil {
		return nil, errInvalidSlug
	}

	wissKI = new(wisski.WissKI)
	instances.use(wissKI)

	wissKI.Liquid.Instance.Slug = slug
	wissKI.Liquid.Instance.FilesystemBase = filepath.Join(instances.Path(), wissKI.Domain())

	wissKI.Liquid.Instance.OwnerEmail = ""
	wissKI.Liquid.Instance.AutoBlindUpdateEnabled = true

	// sql

	wissKI.Liquid.Instance.SqlDatabase = instances.Config.SQL.DataPrefix + slug
	wissKI.Liquid.Instance.SqlUsername = instances.Config.SQL.UserPrefix + slug

	wissKI.Liquid.Instance.SqlPassword, err = instances.Config.NewPassword()
	if err != nil {
		return nil, err
	}

	// triplestore

	wissKI.Liquid.Instance.GraphDBRepository = instances.Config.TS.DataPrefix + slug
	wissKI.Liquid.Instance.GraphDBUsername = instances.Config.TS.UserPrefix + slug

	wissKI.Liquid.Instance.GraphDBPassword, err = instances.Config.NewPassword()
	if err != nil {
		return nil, err
	}

	// drupal

	wissKI.Liquid.DrupalUsername = "admin" // TODO: Change this!

	wissKI.Liquid.DrupalPassword, err = instances.Config.NewPassword()
	if err != nil {
		return nil, err
	}

	// docker image
	wissKI.Liquid.Instance.DockerBaseImage, err = models.GetBaseImage(phpversion)
	if err != nil {
		return nil, err
	}

	// store the instance in the object and return it!
	return wissKI, nil
}

var restrictedSlugs = []string{"www", "admin"}

// IsValidSlug checks if slug represents a valid slug for an instance.
func (instances *Instances) IsValidSlug(slug string) (string, error) {
	// check that it is a slug
	err := validators.ValidateSlug(&slug, "")
	if err != nil {
		return "", errInvalidSlug
	}
	for _, rs := range restrictedSlugs {
		if strings.EqualFold(rs, slug) {
			return "", errRestrictedSlug
		}
	}

	// return the slug
	return slug, nil
}
