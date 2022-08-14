// Package hostname provides hostname
package hostname

import (
	"os"

	"github.com/Showmax/go-fqdn"
)

// FQDN attempts to return the fully qualified domain name of the host system.
// If an error occurs, may fall back to the empty string.
func FQDN() string {

	// try the hostname function
	{
		fqdn, err := fqdn.FqdnHostname()
		if err == nil {
			return fqdn
		}
	}

	// fallback to os hostname
	{
		hostname, err := os.Hostname()
		if err == nil {
			return hostname
		}
	}

	// use the empty string
	return ""
}
