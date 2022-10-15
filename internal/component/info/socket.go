package info

import (
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/component/snapshots"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/tkw1536/goprogram/status"
	"github.com/tkw1536/goprogram/stream"
)

type instanceActionFunc = func(info *Info, instance instances.WissKI, str stream.IOStream) error

var socketInstanceActions = map[string]instanceActionFunc{
	"snapshot": func(info *Info, instance instances.WissKI, str stream.IOStream) error {
		return info.SnapshotManager.MakeExport(
			str,
			snapshots.ExportTask{
				Dest:     "",
				Instance: &instance,

				StagingOnly: false,
			},
		)
	},
	"rebuild": func(_ *Info, instance instances.WissKI, str stream.IOStream) error {
		return instance.Build(str, true)
	},
	"update": func(_ *Info, instance instances.WissKI, str stream.IOStream) error {
		return instance.BlindUpdate(str)
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

func (info *Info) handleInstanceAction(conn httpx.WebSocketConnection, action instanceActionFunc) {

	// read the slug
	slug, ok := <-conn.Read()
	if !ok {
		conn.WriteText("Error reading slug")
		return
	}

	// resolve the instance
	instance, err := info.Instances.WissKI(string(slug.Bytes))
	if err != nil {
		conn.WriteText("Instance not found")
		return
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

	// and perform the action
	{
		err := action(info, instance, str)
		if err != nil {
			str.EPrintln(err)
			return
		}
		str.Println("done")
	}
}
