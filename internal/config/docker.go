package config

type DockerConfig struct {
	// name of docker network to use
	Network string `yaml:"network" default:"distillery" validate:"nonempty"`
}
