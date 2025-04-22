//spellchecker:words config
package config

//spellchecker:words github wisski distillery internal config validators
import "github.com/FAU-CDI/wisski-distillery/internal/config/validators"

// HomeConfig determines options for the homepage of the distillery.
type HomeConfig struct {
	Title        string          `yaml:"title" default:"WissKI Distillery" validate:"nonempty"`
	SelfRedirect *validators.URL `yaml:"redirect" default:"https://github.com/FAU-CDI/wisski-distillery" validate:"https"`
	List         HomeListConfig  `yaml:"list" recurse:"true"`
}

type HomeListConfig struct {
	// Is the list enabled for public visits?
	Public validators.NullableBool `yaml:"public" default:"true" validate:"bool"`
	// Is the list enabled for signed-in visits?
	Private validators.NullableBool `yaml:"private" default:"true" validate:"bool"`
	// Title of the list whenever it is shown
	Title string `yaml:"title" default:"WissKIs on this Distillery" validate:"nonempty"`
}
