package ssh2

//spellchecker:words context github gliderlabs
import (
	"context"
	"io"

	"github.com/gliderlabs/ssh"
)

const (
	etx rune = 3
	eot rune = 4
)

// Server returns an ssh server that implements the main ssh server
func (ssh2 *SSH2) Server(ctx context.Context, privateKeyPath string, progress io.Writer) (*ssh.Server, error) {
	var server ssh.Server

	if err := ssh2.setupHostKeys(progress, ctx, privateKeyPath, &server); err != nil {
		return nil, err
	}

	ssh2.setupForwardHandler(&server)
	ssh2.setupHandler(&server)
	ssh2.setupAuth(&server)

	return &server, nil
}
