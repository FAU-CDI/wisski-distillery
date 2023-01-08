package custom

import (
	_ "embed"
	"html/template"
	"text/template/parse"

	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

const footerName = "footer"

//go:embed "footer.html"
var defaultTemplateStr string
var defaultTemplate = template.Must(template.New("footer.html").Parse(defaultTemplateStr))

// Template creates a copy of template with shared template parts updated accordingly.
// Any template using this should use one of the template contexts in this package.
func (custom *Custom) Template(tpl *template.Template) *template.Template {
	tree := custom.footerTemplate()

	clone := template.Must(tpl.Clone())                 // create a clone of the template
	template.Must(clone.AddParseTree(footerName, tree)) // add the parse tree to it
	return clone                                        // and return the tree
}

// footerTemplate returns a new copy of the footer template
func (custom *Custom) footerTemplate() *parse.Tree {
	footer, err := (func() (*template.Template, error) {
		data, err := environment.ReadFile(custom.Environment, custom.FooterPath())
		if err != nil {
			return nil, err
		}
		return template.New("footer.html").Parse(string(data))
	})()

	if err != nil {
		return defaultTemplate.Tree.Copy()
	}

	return footer.Tree
}
