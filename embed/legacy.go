// Package embed contains embedded resources
package embed

import (
	"embed"
)

// ResourceEmbed contains all the resources required by the WissKI-Distillery package.
//go:embed all:resources
var ResourceEmbed embed.FS
