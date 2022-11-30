package info

import (
	"context"
	"encoding/json"
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

	HandleInteractive func(ctx context.Context, info *Info, instance *wisski.WissKI, out io.Writer, params ...string) error
	HandleResult      func(ctx context.Context, info *Info, instance *wisski.WissKI, params ...string) (value any, err error)
}

var socketInstanceActions = map[string]InstanceAction{
	"snapshot": {
		HandleInteractive: func(ctx context.Context, info *Info, instance *wisski.WissKI, out io.Writer, params ...string) error {
			return info.Exporter.MakeExport(
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
		HandleInteractive: func(ctx context.Context, _ *Info, instance *wisski.WissKI, out io.Writer, params ...string) error {
			return instance.Barrel().Build(ctx, out, true)
		},
	},
	"update": {
		HandleInteractive: func(ctx context.Context, _ *Info, instance *wisski.WissKI, out io.Writer, params ...string) error {
			return instance.Drush().Update(ctx, out)
		},
	},
	"cron": {
		HandleInteractive: func(ctx context.Context, _ *Info, instance *wisski.WissKI, str io.Writer, params ...string) error {
			return instance.Drush().Cron(ctx, str)
		},
	},
}

func (info *Info) serveSocket(conn httpx.WebSocketConnection) {
	// read the next message to act on
	message, ok := <-conn.Read()
	if !ok {
		return
	}

	// perform an action if it exists!
	if action, ok := socketInstanceActions[string(message.Bytes)]; ok {
		info.handleInstanceAction(conn, action)
		return
	}
}

var instanceParamsTimeout = time.Second

func (info *Info) handleInstanceAction(conn httpx.WebSocketConnection, action InstanceAction) {

	// read the slug
	slug, ok := <-conn.Read()
	if !ok {
		<-conn.WriteText("Error reading slug")
		return
	}

	// resolve the instance
	instance, err := info.Instances.WissKI(conn.Context(), string(slug.Bytes))
	if err != nil {
		<-conn.WriteText("Instance not found")
		return
	}

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
		err := action.HandleInteractive(conn.Context(), info, instance, writer, params...)
		if err != nil {
			fmt.Fprintln(writer, err)
			return
		}
		fmt.Fprintln(writer, "done")
	}

	// handle the result computation
	if action.HandleResult != nil {
		result, err := action.HandleResult(conn.Context(), info, instance, params...)
		if err != nil {
			fmt.Fprintln(writer, "false")
			return
		}
		data, err := json.Marshal(result)
		if err != nil {
			fmt.Fprintln(writer, "false")
			return
		}
		fmt.Fprintln(writer, "true")
		fmt.Fprintln(writer, data)
	}

}
