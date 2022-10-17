package wisski

import (
	"time"

	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/stream"
)

var errCronFailed = exit.Error{
	Message:  "Failed to run cron script for instance %q: exited with code %s",
	ExitCode: exit.ExitGeneric,
}

func (wisski *WissKI) Cron(io stream.IOStream) error {
	code, err := wisski.Shell(io, "/runtime/cron.sh")
	if err != nil {
		io.EPrintln(err)
	}
	if code != 0 {
		// keep going, because we want to run as many crons as possible
		err = errBlindUpdateFailed.WithMessageF(wisski.Slug, code)
		io.EPrintln(err)
	}

	return nil
}

func (wisski *WissKI) LastCron(server *PHPServer) (t time.Time, err error) {
	var timestamp int64
	err = wisski.EvalPHPCode(server, &timestamp, `$val = \Drupal::state()->get('system.cron_last'); return $val; `)
	if err != nil {
		return
	}
	return time.Unix(timestamp, 0), nil
}
