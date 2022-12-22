package extras

import (
	"context"
	_ "embed"

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

func (settings *Settings) Set(ctx context.Context, server *phpx.Server, key string, value any) error {
	return settings.Dependencies.PHP.ExecScript(ctx, server, nil, settingsPHP, "set_setting", key, value)
}
