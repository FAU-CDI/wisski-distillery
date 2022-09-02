package env

import "github.com/FAU-CDI/wisski-distillery/internal/stack"

func (dis *Distillery) WebStack() stack.Installable {
	return dis.asCoreStack("web", stack.Installable{
		EnvFileContext: map[string]string{
			"DEFAULT_HOST": dis.Config.DefaultDomain,
		},
	})
}

func (dis *Distillery) WebStackPath() string {
	return dis.WebStack().Dir
}
