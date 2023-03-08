package config

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type ListenConfig struct {
	// Ports are the public addresses to bind to.
	// Each address is automatically multiplexed to serve http, https and ssh traffic.
	// This should typically be port 80 and port 443.
	Ports []uint16 `yaml:"ports" default:"80" validate:"ports"`

	// SSHPort is the port that shows up as the ssh port in various places in the interface.
	// It is automaticalled added to the ports to listen to.
	SSHPort uint16 `yaml:"ssh" default:"80" validate:"port"`
}

// ComposePorts returns a list of ports to be used within a docker-compose.yml file.
// These can be used to forward all ports to the internal port.
func (lc ListenConfig) ComposePorts(internal string) []string {
	// sort and uniquify ports
	ports := append([]uint16{lc.SSHPort}, lc.Ports...)
	slices.Sort(ports)
	ports = slices.Compact(ports)

	forwards := make([]string, len(ports))
	for i, port := range ports {
		forwards[i] = fmt.Sprintf("%d:%s", port, internal)
	}
	return forwards
}
