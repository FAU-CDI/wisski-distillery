package info

import (
	"github.com/FAU-CDI/wisski-distillery/internal/component/snapshots"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/tkw1536/goprogram/status"
	"github.com/tkw1536/goprogram/stream"
)

func (info *Info) serveSocket(conn httpx.WebSocketConnection) {
	// read the next message to act on
	message, ok := <-conn.Read()
	if !ok {
		return
	}

	switch string(message.Bytes) {
	case "snapshot":
		slug, ok := <-conn.Read()
		if !ok {
			return
		}
		info.serverSocketSnapshot(string(slug.Bytes), info.socketWriter(conn))
	case "rebuild":
		slug, ok := <-conn.Read()
		if !ok {
			return
		}
		info.serverSocketRebuild(string(slug.Bytes), info.socketWriter(conn))
	}
}

func (*Info) socketWriter(conn httpx.WebSocketConnection) *status.LineBuffer {
	return &status.LineBuffer{
		Line: func(line string) {
			<-conn.WriteText(line)
		},
		FlushLineOnClose: true,
	}
}

func (info *Info) serverSocketSnapshot(slug string, writer *status.LineBuffer) {
	stream := stream.NewIOStream(writer, writer, nil, 0)

	// get the wisski
	wissKI, err := info.Instances.WissKI(slug)
	if err != nil {
		stream.EPrintln(err)
		return
	}

	{
		err := info.SnapshotManager.MakeExport(
			stream,
			snapshots.ExportTask{
				Dest:     "",
				Instance: &wissKI,

				StagingOnly: false,
			},
		)
		if err != nil {
			stream.EPrintln(err)
			return
		}
	}
	stream.Println("Done")

}

func (info *Info) serverSocketRebuild(slug string, writer *status.LineBuffer) {
	stream := stream.NewIOStream(writer, writer, nil, 0)

	// get the wisski
	wissKI, err := info.Instances.WissKI(slug)
	if err != nil {
		stream.EPrintln(err)
		return
	}

	{
		err := wissKI.Build(stream, true)
		if err != nil {
			stream.EPrintln(err)
			return
		}
	}
	stream.Println("Done")

}
