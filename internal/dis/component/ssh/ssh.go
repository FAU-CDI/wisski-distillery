package ssh

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/pkg/sshx"
	"github.com/gliderlabs/ssh"
	"github.com/tkw1536/goprogram/stream"
	"github.com/tkw1536/proxyssh/feature"
)

type SSH struct {
	component.Base
	Instances *instances.Instances
}

const (
	etx rune = 3
	eot rune = 4
)

const welcomeMessage = `Welcome to the WissKI SSH Server.
You've successfully authenticated, but we don't provide shell access to the main server.
You may use this connection as part of a proxy jump to connect to your server.
For example: 

ssh -J %s:2222 www-data@%s

Press CTRL-C to close this connection.
`

// Server returns an ssh server that implements the main ssh server
func (s *SSH) Server(context context.Context, ios stream.IOStream) (*ssh.Server, error) {
	var server ssh.Server

	banner := fmt.Sprintf(welcomeMessage, s.Config.DefaultDomain, "example."+s.Config.DefaultDomain)

	server.Handle(func(session ssh.Session) {
		io.WriteString(session, banner)

		buffer := bufio.NewReader(session)
		for {
			res, _, err := buffer.ReadRune()
			if err != nil {
				return
			}
			if res == etx || res == eot {
				return
			}
		}

	})
	server.PublicKeyHandler = feature.AuthorizeKeys(
		slogger{IOStream: ios},
		func(ctx ssh.Context) (keys []ssh.PublicKey, err error) {
			keys, err = s.GlobalKeys()
			if err != nil {
				return nil, err
			}

			instances, err := s.Instances.All()
			if err != nil {
				return nil, err
			}

			for _, instance := range instances {
				ikeys, err := instance.SSH().Keys()
				if err != nil {
					continue
				}
				keys = append(keys, ikeys...)
			}
			return keys, nil
		})
	return &server, nil
}

func (s *SSH) GlobalKeys() ([]ssh.PublicKey, error) {
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

type slogger struct {
	stream.IOStream
}

func (s slogger) Print(v ...any) {
	fmt.Fprint(s.Stderr, v...)
}
func (s slogger) Printf(format string, v ...any) {
	fmt.Fprintf(s.Stderr, format, v...)
}
