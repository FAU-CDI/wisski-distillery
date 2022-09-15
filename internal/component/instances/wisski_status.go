package instances

import "github.com/tkw1536/goprogram/stream"

// Info represents some info about this WissKI
type Info struct {
	Slug string // The slug of the instance

	Running bool // is the instance running?
}

// Info returns info about this instance
func (wisski *WissKI) Info() (info Info, err error) {
	info.Slug = wisski.Slug

	ps, err := wisski.Barrel().Ps(stream.FromNil())
	if err != nil {
		return
	}
	info.Running = len(ps) > 0
	return
}
