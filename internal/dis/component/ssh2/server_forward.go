package ssh2

//spellchecker:words github wisski distillery internal component gliderlabs golang crypto gossh
import (
	"fmt"
	"io"
	"net"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

// direct-tcpip data struct as specified in RFC4254, Section 7.2.
type localForwardChannelData struct {
	DestAddr string
	DestPort uint32

	OriginAddr string
	OriginPort uint32
}

// setupForwardHandler sets up the forwarding handler for the ssh server.
func (ssh2 *SSH2) setupForwardHandler(server *ssh.Server) {
	if server.ChannelHandlers == nil {
		server.ChannelHandlers = make(map[string]ssh.ChannelHandler)
		for n, h := range ssh.DefaultChannelHandlers {
			server.ChannelHandlers[n] = h
		}
	}
	server.ChannelHandlers["direct-tcpip"] = ssh2.handleDirectTCP
}

type Intercept struct {
	Description string
	Match       component.HostPort
	Dest        component.HostPort
}

// ExamplePort returns a local port that can be forwarded to without root rights.
func (i Intercept) ExamplePort() uint32 {
	if i.Match.Port < 100 {
		return i.Match.Port * 101
	}
	if i.Match.Port < 1024 {
		return i.Match.Port * 1001
	}
	return i.Match.Port
}
func (i Intercept) Intercept(req component.HostPort) (intercepted bool, ok bool, dest component.HostPort, rejectReason string) {
	if req.Host != i.Match.Host {
		return false, ok, dest, rejectReason
	}

	if req.Port != i.Match.Port {
		return true, false, dest, fmt.Sprintf("%s listens on port %d", i.Description, i.Match.Port)
	}
	return true, true, i.Dest, ""
}

func (ssh2 *SSH2) Intercepts() []Intercept {
	upstream := component.GetStill(ssh2).Upstream
	return ssh2.interceptsC.Get(func() []Intercept {
		return []Intercept{
			{Description: "Triplestore", Match: component.HostPort{Host: "triplestore", Port: 7200}, Dest: upstream.Triplestore},
			{Description: "SQL", Match: component.HostPort{Host: "sql", Port: 3306}, Dest: upstream.SQL},
			{Description: "PHPMyAdmin", Match: component.HostPort{Host: "phpmyadmin", Port: 80}, Dest: component.HostPort{Host: "phpmyadmin", Port: 80}},
		}
	})
}

func (ssh2 *SSH2) getForwardDest(req component.HostPort, ctx ssh.Context) (ok bool, dest component.HostPort, rejectReason string) {
	// check all the intercepts first
	for _, i := range ssh2.Intercepts() {
		intercepted, ok, dest, rejectReason := i.Intercept(req)
		if !intercepted {
			continue
		}
		return ok, dest, rejectReason
	}

	config := component.GetStill(ssh2).Config

	// then check the instances
	slug, ok := config.HTTP.SlugFromHost(req.Host)
	if !ok || req.Port != 22 || !hasPermission(ctx, slug) {
		return false, dest, "permission denied"
	}

	return true, component.HostPort{Host: slug + "." + config.HTTP.PrimaryDomain + ".wisski", Port: 22}, ""
}

// handleDirectTCP handles a direct tcp connection for the server.
//
// #nosec G104
//
//nolint:errcheck // no way to report error
func (ssh2 *SSH2) handleDirectTCP(srv *ssh.Server, conn *gossh.ServerConn, newChan gossh.NewChannel, ctx ssh.Context) {
	d := localForwardChannelData{}
	if err := gossh.Unmarshal(newChan.ExtraData(), &d); err != nil {
		newChan.Reject(gossh.ConnectionFailed, "error parsing forward data: "+err.Error())
		return
	}

	ok, dest, rejectReason := ssh2.getForwardDest(component.HostPort{Host: d.DestAddr, Port: d.DestPort}, ctx)
	if !ok {
		newChan.Reject(gossh.Prohibited, rejectReason)
		return
	}

	var dialer net.Dialer
	dconn, err := dialer.DialContext(ctx, "tcp", dest.String())
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
