package env

import "github.com/FAU-CDI/wisski-distillery/internal/stack"

// SelfComponent represents the 'self' layer belonging to a distillery
type SelfComponent struct {
	dis *Distillery
}

// Self returns the SelfComponent belonging to this distillery
func (dis *Distillery) Self() SelfComponent {
	return SelfComponent{dis: dis}
}

func (SelfComponent) Name() string {
	return "self"
}

func (SelfComponent) Context(parent stack.InstallationContext) stack.InstallationContext {
	return parent
}

func (sc SelfComponent) Stack() stack.Installable {
	TARGET := "https://github.com/FAU-CDI/wisski-distillery"
	if sc.dis.Config.SelfRedirect != nil {
		TARGET = sc.dis.Config.SelfRedirect.String()
	}

	return sc.dis.makeComponentStack(sc, stack.Installable{
		EnvContext: map[string]string{
			"VIRTUAL_HOST":      sc.dis.DefaultVirtualHost(),
			"LETSENCRYPT_HOST":  sc.dis.DefaultLetsencryptHost(),
			"LETSENCRYPT_EMAIL": sc.dis.Config.CertbotEmail,
			"TARGET":            TARGET,
			"OVERRIDES_FILE":    sc.dis.Config.SelfOverridesFile,
		},
	})
}

func (sc SelfComponent) Path() string {
	return sc.Stack().Dir
}
