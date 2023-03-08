package ssh2

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/gliderlabs/ssh"
)

func (ssh2 *SSH2) setupHandler(server *ssh.Server) {
	server.Handle(ssh2.handleConnection)
}

const welcomeMessage = `
__        ___         _  _____   ____  _     _   _ _ _
\ \      / (_)___ ___| |/ /_ _| |  _ \(_)___| |_(_) | | ___ _ __ _   _
 \ \ /\ / /| / __/ __| ' / | |  | | | | / __| __| | | |/ _ \ '__| | | |
  \ V  V / | \__ \__ \ . \ | |  | |_| | \__ \ |_| | | |  __/ |  | |_| |
   \_/\_/  |_|___/___/_|\_\___| |____/|_|___/\__|_|_|_|\___|_|   \__, |
                                                                 |___/

Welcome to the WissKI SSH Server.
You've successfully authenticated, but we don't provide shell access to
the main server. You may use this connection as part of a proxy jump to
connect to your WissKI Instance.

In the following we will provide instructions on how to connect to your
distillery instance via this server. We will assume

${SLUG}

is the name of the WissKI you want to you want to connect to.

From a linux (or mac, or windows 11) command line you may use: 

ssh -J ${DOMAIN}:${PORT} www-data@${HOSTNAME}

You may also place the following into your $HOME/.ssh/config file:

Host *.${DOMAIN}
  ProxyJump ${DOMAIN}.proxy
  User www-data
Host ${DOMAIN}.proxy
  User www-data
  Hostname ${DOMAIN}
  Port ${PORT}

and then connect simply via:

ssh ${HOSTNAME}

On windows you should use the "ssh" executable from the command line if
available. 

If you must, you can also use Putty. 

THIS IS NOT RECOMMENDED AND NOT OFFICIALLY SUPPORTED

First make sure your SSH Key is configured under Connection > Auth > 
Credentials. Then configure a proxy under Connection > Proxy. The Proxy 
Hostname should be

${DOMAIN}

and the port "2222". The proxy type should be "SSH to proxy and use
port forwarding". Then you may enter the hostname 

www-data@${HOSTNAME}

with port 22. 

Press CTRL-C to close this connection.
`

func (ssh2 *SSH2) handleConnection(session ssh.Session) {
	slug, _ := getAnyPermission(session.Context())

	banner := welcomeMessage
	for _, oldnew := range [][2]string{
		{"${SLUG}", slug},
		{"${DOMAIN}", ssh2.Config.HTTP.PrimaryDomain},
		{"${HOSTNAME}", slug + "." + ssh2.Config.HTTP.PrimaryDomain},
		{"${PORT}", strconv.FormatUint(uint64(ssh2.Config.Listen.AdvertisedSSHPort), 10)},
	} {
		banner = strings.ReplaceAll(banner, oldnew[0], oldnew[1])
	}

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
