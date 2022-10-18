package barrel

import "github.com/tkw1536/goprogram/stream"

// Running checks if this WissKI is currently running.
func (barrel *Barrel) Running() (bool, error) {
	ps, err := barrel.Stack().Ps(stream.FromNil())
	if err != nil {
		return false, err
	}
	return len(ps) > 0, nil
}
