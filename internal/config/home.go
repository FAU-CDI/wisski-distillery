//spellchecker:words config
package config

//spellchecker:words github wisski distillery internal config validators
import "github.com/FAU-CDI/wisski-distillery/internal/config/validators"

// HomeConfig determines options for the homepage of the distillery.
type HomeConfig struct {
	Title        string          `default:"WissKI Distillery"                            validate:"nonempty" yaml:"title"`
	SelfRedirect *validators.URL `default:"https://github.com/FAU-CDI/wisski-distillery" validate:"https"    yaml:"redirect"`
	List         HomeListConfig  `recurse:"true"                                         yaml:"list"`
}

type HomeListConfig struct {
	// Is the list enabled for public visits?
	Public validators.NullableBool `default:"true" validate:"bool" yaml:"public"`
	// Is the list enabled for signed-in visits?
	Private validators.NullableBool `default:"true" validate:"bool" yaml:"private"`
	// Title of the list whenever it is shown
	Title string `default:"WissKIs on this Distillery" validate:"nonempty" yaml:"title"`
}
