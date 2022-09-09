package env

import "github.com/FAU-CDI/wisski-distillery/internal/stack"

// SSHComponent represents the 'ssh' layer belonging to a distillery
type SSHComponent struct {
	dis *Distillery
}

// SSH returns the SSHComponent belonging to this distillery
func (dis *Distillery) SSH() SSHComponent {
	return SSHComponent{dis: dis}
}

func (SSHComponent) Name() string {
	return "ssh"
}

func (ssh SSHComponent) Stack() stack.Installable {
	return ssh.dis.makeComponentStack(ssh, stack.Installable{})
}

func (SSHComponent) Context(parent stack.InstallationContext) stack.InstallationContext {
	return parent
}

func (ssh SSHComponent) Path() string {
	return ssh.Stack().Dir
}
