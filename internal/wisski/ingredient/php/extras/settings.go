package extras

import (
	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
)

type Settings struct {
	ingredient.Base

	PHP *php.PHP
}

//go:embed settings.php
var settingsPHP string

func (settings *Settings) Get(server *php.Server, key string) (value any, err error) {
	err = settings.PHP.ExecScript(server, &value, settingsPHP, "get_setting", key)
	return
}

func (settings *Settings) Set(server *php.Server, key string, value any) error {
	return settings.PHP.ExecScript(server, nil, settingsPHP, "set_setting", key, value)
}
