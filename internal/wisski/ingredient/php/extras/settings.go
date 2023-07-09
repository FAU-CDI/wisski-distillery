package extras

import (
	"context"
	_ "embed"
	"errors"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
)

type Settings struct {
	ingredient.Base
	Dependencies struct {
		PHP *php.PHP
	}
}

//go:embed settings.php
var settingsPHP string

func (settings *Settings) Get(ctx context.Context, server *phpx.Server, key string) (value any, err error) {
	err = settings.Dependencies.PHP.ExecScript(ctx, server, &value, settingsPHP, "get_setting", key)
	return
}

var errFailedToSetSetting = errors.New("failed to update setting")

func (settings *Settings) Set(ctx context.Context, server *phpx.Server, key string, value any) error {
	var ok bool
	err := settings.Dependencies.PHP.ExecScript(ctx, server, &ok, settingsPHP, "set_setting", key, value)
	if err == nil && !ok {
		err = errFailedToSetSetting
	}
	return err
}

var errFailedToSetTrustedDomain = errors.New("failed to set trusted domain")

func (settings *Settings) SetTrustedDomain(ctx context.Context, server *phpx.Server, domain string) error {
	var ok bool

	err := settings.Dependencies.PHP.ExecScript(ctx, server, &ok, settingsPHP, "set_trusted_domain", domain)
	if err == nil && !ok {
		err = errFailedToSetTrustedDomain
	}
	return err
}
