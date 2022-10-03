package smartp

import (
	"io"

	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/status"
	"github.com/tkw1536/goprogram/stream"
)

// Run runs f over all items with the given paralllelism.
// When parallel is 1, runs items sequentially with a full input / output stream.
func Run[T any](ios stream.IOStream, parallel int, f func(value T, stream stream.IOStream) error, items []T, opts ...Option[T]) error {

	// create a group
	var group status.Group[T, error]
	group.HandlerLimit = parallel

	// apply all the options
	isParallel := parallel != 1
	for _, opt := range opts {
		group = opt(isParallel, group)
	}

	// setup the default prefix string
	if group.PrefixString == nil {
		group.PrefixString = status.DefaultPrefixString[T]
	}

	// if we are running sequentially
	// then just iterate over the items
	if !isParallel {
		for index, item := range items {
			err := logging.LogOperation(func() error {
				return f(item, ios)
			}, ios, "%v", group.PrefixString(item, index))
			if err != nil {
				return err
			}
		}

		return nil
	}

	// if we are running in parallel, setup a handler
	group.Handler = func(item T, index int, writer io.Writer) error {
		ios := stream.NewIOStream(writer, writer, nil, 0)
		return f(item, ios)
	}

	// create a new status display
	st := status.NewWithCompat(ios.Stdout, 0)
	st.Start()
	defer st.Stop()

	// and use it!
	return status.UseErrorGroup(st, group, items)
}

// Option represents an option of a StatusGroup
type Option[T any] func(bool, status.Group[T, error]) status.Group[T, error]

// SmartMessage returns an option that sets the display of the provided item to the given handler
func SmartMessage[T any](handler func(value T) string) Option[T] {
	return func(p bool, s status.Group[T, error]) status.Group[T, error] {
		s.PrefixString = func(item T, index int) string {
			message := handler(item)
			if p {
				return "[" + message + "]: "
			}
			return message
		}
		s.PrefixAlign = true
		return s
	}
}
