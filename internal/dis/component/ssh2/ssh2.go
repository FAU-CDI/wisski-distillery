package ssh2

import (
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/pkg/sshx"
	"github.com/gliderlabs/ssh"
)

type SSH2 struct {
	component.Base
	Instances *instances.Instances
}

// GlobalKeys returns the global authorized keys
func (s *SSH2) GlobalKeys() ([]ssh.PublicKey, error) {
	file, err := s.Environment.Open(s.Config.GlobalAuthorizedKeysFile)
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return sshx.ParseAllKeys(bytes), nil
}
