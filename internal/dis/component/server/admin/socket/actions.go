package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/provision"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
)

// non-instance specific actions
var actions = map[string]SocketAction{
	"backup": {
		HandleInteractive: func(ctx context.Context, socket *Sockets, out io.Writer, params ...string) error {
			return socket.Dependencies.Exporter.MakeExport(
				ctx,
				out,
				exporter.ExportTask{
					Dest:     "",
					Instance: nil,

					StagingOnly: false,
				},
			)
		},
	},
	"provision": {
		NumParams: 1,
		HandleInteractive: func(ctx context.Context, sockets *Sockets, out io.Writer, params ...string) error {

			// read the flags of the instance to be provisioned
			var flags provision.ProvisionFlags
			if err := json.Unmarshal([]byte(params[0]), &flags); err != nil {
				return err
			}

			instance, err := sockets.Dependencies.Provision.Provision(
				out,
				ctx,
				flags,
			)
			if err != nil {
				return err
			}

			fmt.Fprintf(out, "URL:      %s\n", instance.URL().String())
			fmt.Fprintf(out, "Username: %s\n", instance.DrupalUsername)
			fmt.Fprintf(out, "Password: %s\n", instance.DrupalPassword)

			return nil
		},
	},
}

// socket specific actions
var iActions = map[string]IAction{
	"snapshot": {
		HandleInteractive: func(ctx context.Context, socket *Sockets, instance *wisski.WissKI, out io.Writer, params ...string) error {
			return socket.Dependencies.Exporter.MakeExport(
				ctx,
				out,
				exporter.ExportTask{
					Dest:     "",
					Instance: instance,

					StagingOnly: false,
				},
			)
		},
	},
	"rebuild": {
		HandleInteractive: func(ctx context.Context, _ *Sockets, instance *wisski.WissKI, out io.Writer, params ...string) error {
			return instance.Barrel().Build(ctx, out, true)
		},
	},
	"update": {
		HandleInteractive: func(ctx context.Context, _ *Sockets, instance *wisski.WissKI, out io.Writer, params ...string) error {
			return instance.Drush().Update(ctx, out)
		},
	},
	"cron": {
		HandleInteractive: func(ctx context.Context, _ *Sockets, instance *wisski.WissKI, str io.Writer, params ...string) error {
			return instance.Drush().Cron(ctx, str)
		},
	},
	"start": {
		HandleInteractive: func(ctx context.Context, _ *Sockets, instance *wisski.WissKI, out io.Writer, params ...string) error {
			return instance.Barrel().Stack().Up(ctx, out)
		},
	},
	"stop": {
		HandleInteractive: func(ctx context.Context, _ *Sockets, instance *wisski.WissKI, out io.Writer, params ...string) error {
			return instance.Barrel().Stack().Down(ctx, out)
		},
	},
	"purge": {
		HandleInteractive: func(ctx context.Context, sockets *Sockets, instance *wisski.WissKI, out io.Writer, params ...string) error {
			return sockets.Dependencies.Purger.Purge(ctx, out, instance.Slug)
		},
	},
}

var igActions = func() map[string]SocketAction {
	generics := make(map[string]SocketAction, len(iActions))
	for n, a := range iActions {
		generics[n] = a.AsGenericAction()
	}
	return generics
}()
