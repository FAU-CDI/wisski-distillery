package static

import (
	"embed"
	"encoding/json"
	"html/template"
)

//go:embed "templates/*.html"
var templates embed.FS

var (
	shared *template.Template = template.Must(template.ParseFS(templates, "templates/*.html"))
)

// NewSharedTemplate creates a new template with the given name.
// It will be able to make use of shared templates as well as functions.
func NewSharedTemplate(name string) *template.Template {
	new := template.New(name)
	new.Funcs(template.FuncMap{
		"json": func(data any) (string, error) {
			bytes, err := json.Marshal(data)
			return string(bytes), err
		},
	})
	for _, template := range shared.Templates() {
		new.AddParseTree(template.Tree.Name, template.Tree.Copy())
	}
	return new
}
