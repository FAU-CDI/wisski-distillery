package ssh

import (
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/sshx"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

type SSH struct {
	ingredient.Base
	Barrel *barrel.Barrel
}

var (
	_ ingredient.WissKIFetcher = (*SSH)(nil)
)

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

func (sshx *SSH) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) error {
	if flags.Quick {
		return nil
	}

	keys, err := sshx.Keys()
	if err != nil {
		return err
	}

	info.SSHKeys = make([]string, len(keys))
	for i, key := range keys {
		info.SSHKeys[i] = string(gossh.MarshalAuthorizedKey(key))
	}
	return nil
}
