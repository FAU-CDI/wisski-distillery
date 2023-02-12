package config

import "github.com/FAU-CDI/wisski-distillery/internal/config/validators"

// ThemeConfig determines theming options
type ThemeConfig struct {
	// By default, the default domain redirects to the distillery repository.
	// If you want to change this, set an alternate domain name here.
	SelfRedirect *validators.URL `yaml:"home" default:"https://github.com/FAU-CDI/wisski-distillery" validate:"https"`
}
