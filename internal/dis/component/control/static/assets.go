package static

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
// - Generate new constants of the form {{Name}}
//
// Each asset group should be registered as a parameter to the 'go:generate' line.
type Assets struct {
	Scripts string // <script> tags inserted by the asset
	Styles  string // <link> tags inserted by the asset
}

//go:generate node build.mjs Home User Admin Legal

// MustParse parses a new template from the given source
// and calls [RegisterAssoc] on it.
func (assets *Assets) MustParse(t *template.Template, value string) *template.Template {
	t = template.Must(t.Parse(value))
	assets.RegisterAssoc(t)
	return t
}

// MustParseShared is like [MustParse], but creates a new SharedTemplate instead
func (assets *Assets) MustParseShared(name string, value string) *template.Template {
	return assets.MustParse(NewSharedTemplate(name), value)
}

// RegisterAssoc registers two new associated templates with t.
//
// The template "scripts" will render all script tags required.
// The template "styles" will render all style tags required.
//
// If either template already exists, it will be overwritten.
func (assets *Assets) RegisterAssoc(t *template.Template) {
	t.New("scripts").Parse(assets.Scripts)
	t.New("styles").Parse(assets.Styles)
}
