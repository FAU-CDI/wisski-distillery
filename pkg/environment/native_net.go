package environment

import "net"

func (Native) Listen(network, address string) (net.Listener, error) {
	return net.Listen(network, address)
}

func (Native) Dial(network, address string) (net.Conn, error) {
	return net.Dial(network, address)
}
