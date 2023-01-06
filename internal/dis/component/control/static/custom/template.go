package custom

import (
	_ "embed"
	"html/template"
)

const footerName = "footer"

//go:embed "footer.html"
var defaultTemplateStr string
var defaultTemplate = template.Must(template.New("footer.html").Parse(defaultTemplateStr))

// Template creates a copy of template with shared template parts updated accordingly.
// Any template using this should use one of the template contexts in this package.
func (custom *Custom) Template(tpl *template.Template) *template.Template {
	tree := defaultTemplate.Tree.Copy()

	clone := template.Must(tpl.Clone())                 // create a clone of the template
	template.Must(clone.AddParseTree(footerName, tree)) // add the parse tree to it
	return clone                                        // and return the tree
}
