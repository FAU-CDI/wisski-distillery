package component

import (
	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

type Core struct {
	Environment environment.Environment // environment to use for reading / writing to and from the distillery
	Config      *config.Config          // the configuration of the distillery
}
