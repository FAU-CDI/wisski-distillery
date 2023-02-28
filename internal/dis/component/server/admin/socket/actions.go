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

func (sockets *Sockets) Actions() ActionMap {
	return map[string]Action{
		// generic actions
		"backup": sockets.Generic(0, func(ctx context.Context, sockets *Sockets, in io.Reader, out io.Writer, params ...string) error {
			return sockets.Dependencies.Exporter.MakeExport(
				ctx,
				out,
				exporter.ExportTask{
					Dest:     "",
					Instance: nil,

					StagingOnly: false,
				},
			)
		}),
		"provision": sockets.Generic(1, func(ctx context.Context, sockets *Sockets, in io.Reader, out io.Writer, params ...string) error {
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
		}),

		// instance-specific actions!

		"snapshot": sockets.Instance(0, func(ctx context.Context, socket *Sockets, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) error {
			return socket.Dependencies.Exporter.MakeExport(
				ctx,
				out,
				exporter.ExportTask{
					Dest:     "",
					Instance: instance,

					StagingOnly: false,
				},
			)
		}),
		"rebuild": sockets.Instance(0, func(ctx context.Context, _ *Sockets, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) error {
			return instance.Barrel().Build(ctx, out, true)
		}),
		"update": sockets.Instance(0, func(ctx context.Context, _ *Sockets, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) error {
			return instance.Drush().Update(ctx, out)
		}),
		"cron": sockets.Instance(0, func(ctx context.Context, _ *Sockets, instance *wisski.WissKI, in io.Reader, str io.Writer, params ...string) error {
			return instance.Drush().Cron(ctx, str)
		}),
		"start": sockets.Instance(0, func(ctx context.Context, _ *Sockets, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) error {
			return instance.Barrel().Stack().Up(ctx, out)
		}),
		"stop": sockets.Instance(0, func(ctx context.Context, _ *Sockets, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) error {
			return instance.Barrel().Stack().Down(ctx, out)
		}),
		"purge": sockets.Instance(0, func(ctx context.Context, sockets *Sockets, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) error {
			return sockets.Dependencies.Purger.Purge(ctx, out, instance.Slug)
		}),
	}
}
