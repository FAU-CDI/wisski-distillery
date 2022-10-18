package wisski

import (
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/meta"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/stream"
)

var errBlindUpdateFailed = exit.Error{
	Message:  "Failed to run blind update script for instance %q: exited with code %s",
	ExitCode: exit.ExitGeneric,
}

// BlinUpdate performs a blind update of the given instance
func (wisski *WissKI) BlindUpdate(io stream.IOStream) error {
	code, err := wisski.Shell(io, "/runtime/blind_update.sh")
	if err != nil {
		return errBlindUpdateFailed.WithMessageF(wisski.Slug, environment.ExecCommandError)
	}
	if code != 0 {
		return errBlindUpdateFailed.WithMessageF(wisski.Slug, code)
	}

	return wisski.setLastUpdate()
}

var lastUpdate = meta.StorageFor[int64]("lastUpdate")

func (wisski *WissKI) LastUpdate() (t time.Time, err error) {
	epoch, err := lastUpdate(wisski.storage()).Get()
	if err == meta.ErrMetadatumNotSet {
		return t, nil
	}
	if err != nil {
		return t, err
	}

	// and turn it into time!
	return time.Unix(epoch, 0), nil
}

func (wisski *WissKI) setLastUpdate() error {
	return lastUpdate(wisski.storage()).Set(time.Now().Unix())
}
