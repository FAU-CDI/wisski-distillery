//spellchecker:words manager
package manager

//spellchecker:words maps slices github wisski distillery internal ingredient barrel composer drush system bookkeeping extras
import (
	"maps"
	"slices"

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
	dependencies struct {
		Barrel      *barrel.Barrel
		Bookkeeping *bookkeeping.Bookkeeping

		SystemManager *system.SystemManager

		Composer *composer.Composer
		Drush    *drush.Drush

		Adapters *extras.Adapters
		Settings *extras.Settings
	}
}

// profiles contains the list of default profiles.
var (
	defaultProfile = "Drupal 11"
	profiles       = map[string]Profile{
		"Drupal 9": {
			Description: "Legacy Version of Drupal",

			Drupal: "^9",
			WissKI: "",
			InstallModules: []string{
				"drupal/inline_entity_form:^1.0@RC",
				"drupal/imagemagick",
				"drupal/image_effects",
				"drupal/colorbox",
			},
			EnableModules: []string{
				"drupal/devel:^4.1",
				"drupal/geofield:^1.40",
				"drupal/geofield_map:^2.85",
				"drupal/imce:^2.4",
				"drupal/remove_generator:^2.0",
			},
		},
		"Drupal 10": {
			Description: "Legacy Version Of Drupal",

			Drupal: "^10",
			WissKI: "",
			InstallModules: []string{
				"drupal/inline_entity_form:^1.0@RC",
				"drupal/imagemagick",
				"drupal/image_effects",
				"drupal/colorbox",
			},
			EnableModules: []string{
				"drupal/devel:^5.0",
				"drupal/geofield:^1.56",
				"drupal/geofield_map:^3.0",
				"drupal/imce:^3.0",
				"drupal/remove_generator:^2.0",
			},
		},
		"Drupal 11": {
			Description: "Current Version of Drupal with default packages",

			Drupal: "^11",
			WissKI: "",
			InstallModules: []string{
				"drupal/inline_entity_form:^3.0@RC",
				"drupal/imagemagick",
				"drupal/image_effects",
				"drupal/colorbox",
			},
			EnableModules: []string{
				"drupal/devel:^5.3",
				"drupal/geofield:^1.64",
				"drupal/geofield_map:^11.0",
				"drupal/imce:^3.1",
				"drupal/remove_generator:^2.1",
			},
		},
	}
)

// TODO: All of these should move to the config

func LoadDefaultProfile() Profile {
	return LoadProfile(DefaultProfile())
}

func Profiles() map[string]Profile {
	return maps.Clone(profiles)
}

func LoadProfile(name string) Profile {
	return profiles[name]
}

func HasProfile(name string) bool {
	_, ok := profiles[name]
	return ok
}

func DefaultProfile() string {
	return defaultProfile
}

// Profile represents a profile applied to a WissKI instance of the Distillery.
type Profile struct {
	// Description is a human-readable description for this profile.
	// It is only used by the frontend.
	Description string

	Drupal string // Version of Drupal to use
	WissKI string // Version of WissKI to use

	InstallModules []string // Modules to be installed (but not neccessarily enabled)
	EnableModules  []string // Modules to be installed and enabled
}

// Apply copies over defaults from the other profile to this one.
// If a field is already set, no defaults are copied.
func (profile *Profile) Apply(other Profile) {
	if profile.Drupal == "" {
		profile.Drupal = other.Drupal
	}
	if profile.WissKI == "" {
		profile.WissKI = other.WissKI
	}
	if profile.InstallModules == nil {
		profile.InstallModules = slices.Clone(other.InstallModules)
	}
	if profile.EnableModules == nil {
		profile.EnableModules = slices.Clone(profile.EnableModules)
	}
}

// ApplyDefaults loads some set of defaults.
// If all fields are set, no defaults are applied.
func (profile *Profile) ApplyDefaults() {
	profile.Apply(profiles[defaultProfile])
}
