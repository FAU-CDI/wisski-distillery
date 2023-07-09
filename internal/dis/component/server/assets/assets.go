package assets

import (
	"html/template"
)

// Assets represents a group of assets to be included inside a template.
//
// Assets are generated using the 'build.mjs' script.
// The script is called using 'go:generate', which stores variables in the form of 'Assets{{Name}}' inside this package.
//
// The build script roughly works as follows:
// - Delete any previously generated distribution directory.
// - Bundle the entrypoint sources under 'src/entry/{{Name}}/index.{ts,css}' together with the base './src/base/index.{ts,css}'
// - Store the output inside the 'dist' directory
// - Generate new constants of the form Assets{{Name}}
//
// Each asset group should be registered as a parameter to the 'go:generate' line.
type Assets struct {
	Scripts template.HTML // <script> tags inserted by the asset
	Styles  template.HTML // <link> tags inserted by the asset
}

//go:generate node build.mjs Default User Admin AdminProvision AdminRebuild
