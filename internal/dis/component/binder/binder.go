package binder

import (
	"bytes"
	"io"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/tkw1536/pkglib/yamlx"
	"gopkg.in/yaml.v3"

	_ "embed"
)

type Binder struct {
	component.Base
}

var (
	_ component.Installable = (*Binder)(nil)
)

func (binder *Binder) Path() string {
	return filepath.Join(binder.Still.Config.Paths.Root, "core", "binder")
}

func (binder *Binder) Context(parent component.InstallationContext) component.InstallationContext {
	return parent
}

//go:embed docker-compose.yml
var composeTemplate []byte

func (binder *Binder) buildYML() ([]byte, error) {
	var dockerCompose yaml.Node
	if err := yaml.Unmarshal(composeTemplate, &dockerCompose); err != nil {
		return nil, err
	}

	for dockerCompose.Kind == yaml.DocumentNode {
		dockerCompose = *dockerCompose.Content[0]
	}

	{
		ports := binder.Config.Listen.ComposePorts("8000")
		portsNode, err := yamlx.Marshal(ports)
		if err != nil {
			return nil, err
		}
		if err := yamlx.Replace(&dockerCompose, *portsNode, "services", "binder", "ports"); err != nil {
			return nil, err
		}
	}

	{
		command := binder.Config.HTTP.TCPMuxCommand("0.0.0.0:8000", "http:80", "http:443", "ssh:2222")
		commandNode, err := yamlx.Marshal(command)
		if err != nil {
			return nil, err
		}
		if err := yamlx.Replace(&dockerCompose, *commandNode, "services", "binder", "command"); err != nil {
			return nil, err
		}
	}

	// do the final marshal
	return yaml.Marshal(dockerCompose)
}

func (binder *Binder) Stack() component.StackWithResources {
	return component.MakeStack(binder, component.StackWithResources{
		ReadComposeFile: func() (io.Reader, error) {
			data, err := binder.buildYML()
			if err != nil {
				return nil, err
			}
			return bytes.NewReader(data), nil
		},

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": binder.Config.Docker.Network(),
		},
	})
}
