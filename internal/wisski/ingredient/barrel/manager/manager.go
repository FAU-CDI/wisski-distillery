package manager

import (
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/composer"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/drush"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/system"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/bookkeeping"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php/extras"
)

// Manager manages a profile applied to specific WissKI instances.
type Manager struct {
	ingredient.Base
	Dependencies struct {
		Barrel      *barrel.Barrel
		Bookkeeping *bookkeeping.Bookkeeping

		SystemManager *system.SystemManager

		Composer *composer.Composer
		Drush    *drush.Drush

		Adapters *extras.Adapters
		Settings *extras.Settings
	}
}

// Profile represents a profile applied to a WissKI instance of the Distillery.
type Profile struct {
	Drupal string // Version of Drupal to use
	WissKI string // Version of WissKI to use

	InstallModules []string // Modules to be installed (but not neccessarily enabled)
	EnableModules  []string // Modules to be installed and enabled
}

// DefaultDrupalVersion is the default drupal version
const DefaultDrupalVersion = "^9.0.0"

// ApplyDefaults applies the default settings to missing profile settings.
func (profile *Profile) ApplyDefaults() {
	if profile.Drupal == "" {
		profile.Drupal = DefaultDrupalVersion
	}
	if profile.InstallModules == nil {
		profile.InstallModules = []string{
			"drupal/inline_entity_form:^1.0@RC",
			"drupal/imagemagick",
			"drupal/image_effects",
			"drupal/colorbox",
		}
	}
	if profile.EnableModules == nil {
		profile.EnableModules = []string{
			"drupal/devel:^4.1",
			"drupal/geofield:^1.40",
			"drupal/geofield_map:^2.85",
			"drupal/imce:^2.4",
			"drupal/remove_generator:^2.0",
		}
	}
}
