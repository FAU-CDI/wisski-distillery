package environment

import (
	"context"
	"net"
)

func (*Native) Listen(network, address string) (net.Listener, error) {
	return net.Listen(network, address)
}

func (*Native) DialContext(context context.Context, network, address string) (net.Conn, error) {
	var d net.Dialer
	return d.DialContext(context, network, address)
}
