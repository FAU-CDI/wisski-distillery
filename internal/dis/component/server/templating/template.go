package templating

import (
	_ "embed"
	"html/template"
	"os"
)

//go:embed "src/footer.html"
var footerHTML string
var footerTemplate = template.Must(template.New("footer.html").Parse(footerHTML))

// GetCustomizable returns either a clone of dflt, or the overriden template with the same name.
func (tpl *Templating) GetCustomizable(dflt *template.Template) *template.Template {
	name := dflt.Name()

	custom, err := (func() (*template.Template, error) {
		data, err := os.ReadFile(tpl.CustomAssetPath(name))
		if err != nil {
			return nil, err
		}
		return template.New(name).Parse(string(data))
	})()
	if err != nil {
		return template.Must(dflt.Clone())
	}
	return custom
}
