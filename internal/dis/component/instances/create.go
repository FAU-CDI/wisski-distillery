//spellchecker:words instances
package instances

//spellchecker:words errors path filepath strings github wisski distillery internal config validators component models
import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/config/validators"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
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
func (instances *Instances) Create(slug string, system models.System) (wissKI *wisski.WissKI, err error) {
	// make sure that the slug is valid!
	slug, err = instances.IsValidSlug(slug)
	if err != nil {
		return nil, errInvalidSlug
	}

	wissKI = new(wisski.WissKI)
	instances.use(wissKI)

	wissKI.Slug = slug
	wissKI.FilesystemBase = filepath.Join(instances.Path(), wissKI.Domain())

	wissKI.OwnerEmail = ""
	wissKI.AutoBlindUpdateEnabled = true

	config := component.GetStill(instances).Config

	// sql

	wissKI.SqlDatabase = config.SQL.DataPrefix + slug
	wissKI.SqlUsername = config.SQL.UserPrefix + slug

	wissKI.SqlPassword, err = config.NewPassword()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new password: %w", err)
	}

	// triplestore

	wissKI.GraphDBRepository = config.TS.DataPrefix + slug
	wissKI.GraphDBUsername = config.TS.UserPrefix + slug

	wissKI.GraphDBPassword, err = config.NewPassword()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new password: %w", err)
	}

	// drupal

	wissKI.DrupalUsername = "admin" // TODO: Change this!

	wissKI.DrupalPassword, err = config.NewPassword()
	if err != nil {
		return nil, fmt.Errorf("failed to get net password: %w", err)
	}

	// docker image
	wissKI.System = system

	// store the instance in the object and return it!
	return wissKI, nil
}

// IsValidSlug checks if slug represents a valid slug for an instance.
func (instances *Instances) IsValidSlug(slug string) (string, error) {
	// check that it is a slug
	err := validators.ValidateSlug(&slug, "")
	if err != nil {
		return "", errInvalidSlug
	}
	for _, rs := range config.RestrictedSlugs {
		if strings.EqualFold(rs, slug) {
			return "", errRestrictedSlug
		}
	}

	// return the slug
	return slug, nil
}
