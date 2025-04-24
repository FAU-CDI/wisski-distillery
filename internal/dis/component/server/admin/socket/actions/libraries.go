//spellchecker:words actions
package actions

//spellchecker:words context github wisski distillery internal component auth scopes pkglib stream
import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/tkw1536/pkglib/stream"
)

// installing additional javascript libraries

type InstallColorboxJS struct {
	component.Base
}

var (
	_ WebsocketInstanceAction = (*InstallColorboxJS)(nil)
)

func (*InstallColorboxJS) Action() InstanceAction {
	return InstanceAction{
		Action: Action{
			Name:      "install-colorbox-js",
			Scope:     scopes.ScopeUserAdmin,
			NumParams: 0,
		},
	}
}

func (*InstallColorboxJS) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) (any, error) {
	if err := instance.Barrel().Shell(ctx, stream.NewIOStream(out, out, nil), "/runtime/install_colorbox.sh"); err != nil {
		return nil, fmt.Errorf("failed to install colorbox: %w", err)
	}
	return nil, nil
}

type InstallDompurifyJS struct {
	component.Base
}

var (
	_ WebsocketInstanceAction = (*InstallDompurifyJS)(nil)
)

func (*InstallDompurifyJS) Action() InstanceAction {
	return InstanceAction{
		Action: Action{
			Name:      "install-dompurify-js",
			Scope:     scopes.ScopeUserAdmin,
			NumParams: 0,
		},
	}
}

func (*InstallDompurifyJS) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) (any, error) {
	if err := instance.Barrel().Shell(ctx, stream.NewIOStream(out, out, nil), "/runtime/install_dompurify.sh"); err != nil {
		return nil, fmt.Errorf("failed to install dompurify: %w", err)
	}
	return nil, nil
}
