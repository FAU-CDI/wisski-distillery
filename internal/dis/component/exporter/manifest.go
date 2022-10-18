package exporter

import (
	"path/filepath"
)

type WithManifest struct {
	Manifest []string
}

func (wm *WithManifest) handleManifest(dest string) (chan<- string, func()) {
	manifest := make(chan string)
	done := make(chan struct{})
	go func() {
		defer close(done)

		for file := range manifest {
			// get the relative path to the root of the manifest.
			// nothing *should* go wrong, but in case it does, use the original path.
			path, err := filepath.Rel(dest, file)
			if err != nil {
				path = file
			}

			// add the manifest
			wm.Manifest = append(wm.Manifest, path)
		}
	}()
	return manifest, func() {
		close(manifest)
		<-done
	}
}
