package ssh2

import (
	"io"
	"net"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

// direct-tcpip data struct as specified in RFC4254, Section 7.2
type localForwardChannelData struct {
	DestAddr string
	DestPort uint32

	OriginAddr string
	OriginPort uint32
}

// seetupForwardHandler sets up the forwarding handler for the ssh server
func (ssh2 *SSH2) setupForwardHandler(server *ssh.Server) {
	if server.ChannelHandlers == nil {
		server.ChannelHandlers = make(map[string]ssh.ChannelHandler)
		for n, h := range ssh.DefaultChannelHandlers {
			server.ChannelHandlers[n] = h
		}
	}
	server.ChannelHandlers["direct-tcpip"] = ssh2.handleDirectTCP
}

// handleDirectTCP handles a direct tcp connection for the server
func (ssh2 *SSH2) handleDirectTCP(srv *ssh.Server, conn *gossh.ServerConn, newChan gossh.NewChannel, ctx ssh.Context) {
	d := localForwardChannelData{}
	if err := gossh.Unmarshal(newChan.ExtraData(), &d); err != nil {
		newChan.Reject(gossh.ConnectionFailed, "error parsing forward data: "+err.Error())
		return
	}

	slug, ok := ssh2.Config.SlugFromHost(d.DestAddr)
	if !ok || d.DestPort != 22 || !hasPermission(ctx, slug) {
		newChan.Reject(gossh.Prohibited, "permission denied")
		return
	}

	// TODO: move this into an instance function somewhere
	dest := net.JoinHostPort(slug+"."+ssh2.Config.DefaultDomain+".wisski", "22")

	var dialer net.Dialer
	dconn, err := dialer.DialContext(ctx, "tcp", dest)
	if err != nil {
		newChan.Reject(gossh.ConnectionFailed, err.Error())
		return
	}

	ch, reqs, err := newChan.Accept()
	if err != nil {
		dconn.Close()
		return
	}
	go gossh.DiscardRequests(reqs)

	go func() {
		defer ch.Close()
		defer dconn.Close()
		io.Copy(ch, dconn)
	}()
	go func() {
		defer ch.Close()
		defer dconn.Close()
		io.Copy(dconn, ch)
	}()
}
