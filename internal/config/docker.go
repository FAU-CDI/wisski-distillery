//spellchecker:words config
package config

type DockerConfig struct {
	NetworkPrefix string `default:"distillery" validate:"nonempty" yaml:"network"`
}

// Networks returns a list of all docker networks to be created for purposes of the distillery.
func (dc DockerConfig) Networks() []string {
	return []string{dc.Network()}
}

// Network returns the name of the default network to attach all docker containers to.
func (dc DockerConfig) Network() string {
	return dc.NetworkPrefix
}
