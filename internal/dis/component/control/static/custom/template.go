package custom

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
func (custom *Custom) Template(tpl *template.Template) *template.Template {
	// create a clone of the template
	clone := template.Must(tpl.Clone())

	// add all the fixed parse trees
	footerTree := custom.getTemplateAsset(defaultFooterTemplate)
	template.Must(clone.AddParseTree(footerName, footerTree))

	// optionally add the about asset
	if aboutTree := custom.readTemplateAsset("about.html"); clone.Lookup(aboutName) != nil && aboutTree != nil {
		template.Must(clone.AddParseTree(aboutName, aboutTree))
	}
	return clone // and return the tree
}

// getTemplateAsset returns an overridable template asset.
//
// If the asset named can successfully be parsed, it is returned.
// If it can not be parsed, the default template is returned.
func (custom *Custom) getTemplateAsset(dflt *template.Template) *parse.Tree {
	tree := custom.readTemplateAsset(dflt.Name())
	if tree == nil {
		return dflt.Tree.Copy()
	}
	return tree
}

// readTemplateAsset is like getTemplateAssets, but takes an explicit name to read.
// when the asset does not exist, or cannot be opened, returns nil.
func (custom *Custom) readTemplateAsset(name string) *parse.Tree {
	template, err := (func() (*template.Template, error) {
		data, err := environment.ReadFile(custom.Environment, custom.CustomAssetPath(name))
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
