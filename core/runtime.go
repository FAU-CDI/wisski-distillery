package core

import (
	"embed"
	_ "embed"
)

// Runtime contains runtime resources to be installed into any instance
//go:embed all:runtime
var Runtime embed.FS
