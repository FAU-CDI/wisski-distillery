package ssh

import (
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/sshx"
	"github.com/gliderlabs/ssh"
)

type SSH struct {
	ingredient.Base
	Barrel *barrel.Barrel
}

func (ssh *SSH) Keys() ([]ssh.PublicKey, error) {
	file, err := ssh.Environment.Open(ssh.Barrel.AuthorizedKeysPath())
	if environment.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return sshx.ParseAllKeys(bytes), nil
}
