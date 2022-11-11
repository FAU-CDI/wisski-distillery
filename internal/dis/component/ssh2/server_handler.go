package ssh2

import (
	"bufio"
	"fmt"
	"io"

	"github.com/gliderlabs/ssh"
)

func (ssh2 *SSH2) setupHandler(server *ssh.Server) {
	server.Handle(ssh2.handleConnection)
}

func (ssh2 *SSH2) handleConnection(session ssh.Session) {
	slug, _ := getAnyPermission(session.Context())
	banner := fmt.Sprintf(welcomeMessage, ssh2.Config.DefaultDomain, slug+"."+ssh2.Config.DefaultDomain)

	io.WriteString(session, banner)

	// wait until the user closes
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
}
