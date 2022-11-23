package solr

import (
	"embed"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

type Keycloak struct {
	component.Base

	SQL *sql.SQL
}

func (k *Keycloak) Path() string {
	return filepath.Join(k.Still.Config.DeployRoot, "core", "keycloak")
}

func (*Keycloak) Context(parent component.InstallationContext) component.InstallationContext {
	return parent
}

//go:embed all:keycloak
//go:embed keycloak.env
var resources embed.FS

func (kc *Keycloak) Stack(env environment.Environment) component.StackWithResources {
	return component.MakeStack(kc, env, component.StackWithResources{
		Resources:   resources,
		ContextPath: "keycloak",

		EnvPath: "keycloak.env",
		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME":     kc.Config.DockerNetworkName,
			"KEYCLOAK_ADMIN":          kc.Config.KeycloakAdminUser,
			"KEYCLOAK_ADMIN_PASSWORD": kc.Config.KeycloakAdminPassword,
		},

		MakeDirs: []string{
			filepath.Join("data"),
		},
	})
}
