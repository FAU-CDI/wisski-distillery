package env

import "github.com/FAU-CDI/wisski-distillery/internal/stack"

func (dis *Distillery) SSHStack() stack.Installable {
	// TODO: Ensure that .env is copied if needed
	return dis.asCoreStack("ssh", stack.Installable{})
}

func (dis *Distillery) SSHStackPath() string {
	return dis.SSHStack().Dir
}
