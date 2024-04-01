package extras

import (
	"context"
	_ "embed"
	"errors"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
)

type Settings struct {
	ingredient.Base
	dependencies struct {
		PHP *php.PHP
	}
}

//go:embed settings.php
var settingsPHP string

func (settings *Settings) Get(ctx context.Context, server *phpx.Server, key string) (value any, err error) {
	err = settings.dependencies.PHP.ExecScript(ctx, server, &value, settingsPHP, "get_setting", key)
	return
}

var errFailedToSetSetting = errors.New("failed to update setting")

func (settings *Settings) Set(ctx context.Context, server *phpx.Server, key string, value any) error {
	var ok bool
	err := settings.dependencies.PHP.ExecScript(ctx, server, &ok, settingsPHP, "set_setting", key, value)
	if err == nil && !ok {
		err = errFailedToSetSetting
	}
	return err
}

var (
	errFailedToSetTrustedDomain        = errors.New("failed to set trusted domain")
	errFailedInstallDistillerySettings = errors.New("failed to install distillery settings")
)

// SetTrustedDomain configures the trusted domain setting for the given instance.
// Note that this removes any installed distillery settings.
func (settings *Settings) SetTrustedDomain(ctx context.Context, server *phpx.Server, domain string) error {
	var ok bool

	err := settings.dependencies.PHP.ExecScript(ctx, server, &ok, settingsPHP, "set_trusted_domain", domain)
	if err == nil && !ok {
		err = errFailedToSetTrustedDomain
	}
	return err
}

func (settings *Settings) InstallDistillerySettings(ctx context.Context, server *phpx.Server) error {
	var ok bool

	err := settings.dependencies.PHP.ExecScript(ctx, server, &ok, settingsPHP, "install_settings_include", []string{
		barrel.LocalSettingsPath,
		barrel.GlobalSettingsPath,
	})
	if err == nil && !ok {
		err = errFailedInstallDistillerySettings
	}
	return err
}
