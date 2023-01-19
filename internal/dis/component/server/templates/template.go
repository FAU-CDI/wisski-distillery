package templates

import (
	_ "embed"
	"html/template"
	"text/template/parse"

	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

const (
	footerName = "@custom/footer"
	aboutName  = "@custom/about"
)

//go:embed "footer.html"
var footerTemplateStr string
var defaultFooterTemplate = template.Must(template.New("footer.html").Parse(footerTemplateStr))

// Template creates a copy of template with shared template parts updated accordingly.
// Any template using this should use one of the template contexts in this package.
func (tpl *Templating) Template(t *template.Template) *template.Template {
	// TODO: This should not be used!

	// create a clone of the template
	clone := template.Must(t.Clone())

	// add all the fixed parse trees
	footerTree := tpl.getTemplateAsset(defaultFooterTemplate)
	template.Must(clone.AddParseTree(footerName, footerTree))

	// optionally add the about asset
	if aboutTree := tpl.readTemplateAsset("about.html"); clone.Lookup(aboutName) != nil && aboutTree != nil {
		template.Must(clone.AddParseTree(aboutName, aboutTree))
	}
	return clone // and return the tree
}

// getTemplateAsset returns an overridable template asset.
//
// If the asset named can successfully be parsed, it is returned.
// If it can not be parsed, the default template is returned.
func (tpl *Templating) getTemplateAsset(dflt *template.Template) *parse.Tree {
	tree := tpl.readTemplateAsset(dflt.Name())
	if tree == nil {
		return dflt.Tree.Copy()
	}
	return tree
}

// readTemplateAsset is like getTemplateAssets, but takes an explicit name to read.
// when the asset does not exist, or cannot be opened, returns nil.
func (tpl *Templating) readTemplateAsset(name string) *parse.Tree {
	template, err := (func() (*template.Template, error) {
		data, err := environment.ReadFile(tpl.Environment, tpl.CustomAssetPath(name))
		if err != nil {
			return nil, err
		}
		return template.New(name).Parse(string(data))
	})()
	if err != nil {
		return nil
	}
	return template.Tree
}
