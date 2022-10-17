package wisski

import (
	_ "embed"
)

//go:embed php/settings.php
var settingsPHP string

func (wisski *WissKI) GetSettingsPHP(server *PHPServer, key string) (value any, err error) {
	err = wisski.ExecPHPScript(server, &value, settingsPHP, "get_setting", key)
	return
}

func (wisski *WissKI) SetSettingsPHP(server *PHPServer, key string, value any) error {
	return wisski.ExecPHPScript(server, nil, settingsPHP, "set_setting", key, value)
}
