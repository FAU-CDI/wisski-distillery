//spellchecker:words templating
package templating

//spellchecker:words embed html template
import (
	_ "embed"
	"fmt"
	"html/template"
	"os"
)

//go:embed "src/footer.html"
var footerHTML string
var footerTemplate = template.Must(template.New("footer.html").Parse(footerHTML))

// GetCustomizable returns either a clone of dflt, or the overriden template with the same name.
func (tpl *Templating) GetCustomizable(dflt *template.Template) *template.Template {
	custom, err := tpl.getCustomizableTemplate(dflt.Name())
	if err != nil {
		return template.Must(dflt.Clone())
	}
	return custom
}

func (tpl *Templating) getCustomizableTemplate(name string) (*template.Template, error) {
	data, err := os.ReadFile(tpl.CustomAssetPath(name))
	if err != nil {
		return nil, fmt.Errorf("failed to read custom asset: %w", err)
	}
	template, err := template.New(name).Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}
	return template, nil
}
