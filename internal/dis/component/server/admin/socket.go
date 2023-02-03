package admin

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/tkw1536/goprogram/status"
)

type InstanceAction struct {
	NumParams int

	HandleInteractive func(ctx context.Context, info *Admin, instance *wisski.WissKI, out io.Writer, params ...string) error
}

func (ia *InstanceAction) AsGenericAction() GenericAction {
	return GenericAction{
		NumParams: ia.NumParams + 1,
		HandleInteractive: func(ctx context.Context, info *Admin, out io.Writer, params ...string) error {
			instance, err := info.Dependencies.Instances.WissKI(ctx, params[0])
			if err != nil {
				return err
			}

			return ia.HandleInteractive(ctx, info, instance, out, params[1:]...)
		},
	}
}

type GenericAction struct {
	NumParams int

	HandleInteractive func(ctx context.Context, info *Admin, out io.Writer, params ...string) error
}

// non-instance specific actions
var genericActions = map[string]GenericAction{}

// socket specific actions
var socketInstanceActions = map[string]InstanceAction{
	"snapshot": {
		HandleInteractive: func(ctx context.Context, admin *Admin, instance *wisski.WissKI, out io.Writer, params ...string) error {
			return admin.Dependencies.Exporter.MakeExport(
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
		HandleInteractive: func(ctx context.Context, _ *Admin, instance *wisski.WissKI, out io.Writer, params ...string) error {
			return instance.Barrel().Build(ctx, out, true)
		},
	},
	"update": {
		HandleInteractive: func(ctx context.Context, _ *Admin, instance *wisski.WissKI, out io.Writer, params ...string) error {
			return instance.Drush().Update(ctx, out)
		},
	},
	"cron": {
		HandleInteractive: func(ctx context.Context, _ *Admin, instance *wisski.WissKI, str io.Writer, params ...string) error {
			return instance.Drush().Cron(ctx, str)
		},
	},
	"start": {
		HandleInteractive: func(ctx context.Context, _ *Admin, instance *wisski.WissKI, out io.Writer, params ...string) error {
			return instance.Barrel().Stack().Up(ctx, out)
		},
	},
	"stop": {
		HandleInteractive: func(ctx context.Context, _ *Admin, instance *wisski.WissKI, out io.Writer, params ...string) error {
			return instance.Barrel().Stack().Down(ctx, out)
		},
	},
	"purge": {
		HandleInteractive: func(ctx context.Context, admin *Admin, instance *wisski.WissKI, out io.Writer, params ...string) error {
			return admin.Dependencies.Purger.Purge(ctx, out, instance.Slug)
		},
	},
}

var socketGenericActions = func() map[string]GenericAction {
	generics := make(map[string]GenericAction, len(socketInstanceActions))
	for n, a := range socketInstanceActions {
		generics[n] = a.AsGenericAction()
	}
	return generics
}()

func (admin *Admin) serveSocket(conn httpx.WebSocketConnection) {
	// read the next message to act on
	message, ok := <-conn.Read()
	if !ok {
		return
	}

	name := string(message.Bytes)

	// perform a generic action first
	if action, ok := genericActions[name]; ok {
		admin.handleGenericAction(conn, action)
		return
	}

	// then do the socket actions
	if action, ok := socketGenericActions[name]; ok {
		admin.handleGenericAction(conn, action)
	}
}

var instanceParamsTimeout = time.Second

func (admin *Admin) handleGenericAction(conn httpx.WebSocketConnection, action GenericAction) {
	// read the parameters
	params := make([]string, action.NumParams)
	for i := range params {
		select {
		case message, ok := <-conn.Read():
			if !ok {
				<-conn.WriteText("Insufficient parameters")
				return
			}
			params[i] = string(message.Bytes)
		case <-time.After(instanceParamsTimeout):
			<-conn.WriteText("Timed out reading parameters")
			return
		}
	}

	// build a stream
	writer := &status.LineBuffer{
		Line: func(line string) {
			<-conn.WriteText(line)
		},
		FlushLineOnClose: true,
	}
	defer writer.Close()

	// handle the interactive action
	if action.HandleInteractive != nil {
		err := action.HandleInteractive(conn.Context(), admin, writer, params...)
		if err != nil {
			fmt.Fprintln(writer, err)
			return
		}
		fmt.Fprintln(writer, "done")
	}
}
