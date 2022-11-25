package info

import (
	"encoding/json"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/tkw1536/goprogram/status"
	"github.com/tkw1536/goprogram/stream"
)

type InstanceAction struct {
	NumParams int

	HandleInteractive func(info *Info, instance *wisski.WissKI, str stream.IOStream, params ...string) error
	HandleResult      func(info *Info, instance *wisski.WissKI, params ...string) (value any, err error)
}

var socketInstanceActions = map[string]InstanceAction{
	"snapshot": {
		HandleInteractive: func(info *Info, instance *wisski.WissKI, str stream.IOStream, params ...string) error {
			return info.Exporter.MakeExport(
				str,
				exporter.ExportTask{
					Dest:     "",
					Instance: instance,

					StagingOnly: false,
				},
			)
		},
	},
	"rebuild": {
		HandleInteractive: func(_ *Info, instance *wisski.WissKI, str stream.IOStream, params ...string) error {
			return instance.Barrel().Build(str, true)
		},
	},
	"update": {
		HandleInteractive: func(_ *Info, instance *wisski.WissKI, str stream.IOStream, params ...string) error {
			return instance.Drush().Update(str)
		},
	},
	"cron": {
		HandleInteractive: func(_ *Info, instance *wisski.WissKI, str stream.IOStream, params ...string) error {
			return instance.Drush().Cron(str)
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
	instance, err := info.Instances.WissKI(string(slug.Bytes))
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

	str := stream.NewIOStream(writer, writer, nil, 0)

	// handle the interactive action
	if action.HandleInteractive != nil {
		err := action.HandleInteractive(info, instance, str, params...)
		if err != nil {
			str.EPrintln(err)
			return
		}
		str.Println("done")
	}

	// handle the result computation
	if action.HandleResult != nil {
		result, err := action.HandleResult(info, instance, params...)
		if err != nil {
			str.Println("false")
			return
		}
		data, err := json.Marshal(result)
		if err != nil {
			str.Println("false")
			return
		}
		data = append(data, "\n"...)
		str.Println("true")
		str.Stdout.Write(data)
	}

}
