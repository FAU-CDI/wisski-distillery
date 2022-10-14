package static

import (
	"bytes"

	"github.com/FAU-CDI/wisski-distillery/pkg/resources"
)

var EntryHome = mustParseResources("dist/home/index.html")
var EntryControlIndex = mustParseResources("dist/control/index.html")
var EntryControlInstance = mustParseResources("dist/control/instance.html")

// mustParseResources loads the resources from the provided files or panic()s
func mustParseResources(path string) resources.Resources {
	data, err := distStaticFS.ReadFile(path)
	if err != nil {
		panic("mustParseResources: Unable to open " + path)
	}
	return resources.Parse(bytes.NewReader(data))
}
