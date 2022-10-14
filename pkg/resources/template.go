package resources

import (
	"html/template"
	"strings"
)

// MustParse parses a new "html/template" from the given value, and registers the given functions with it.
// When something goes wrong, calls panic()
func (resources *Resources) MustParse(value string) *template.Template {
	return template.Must(resources.RegisterFuncs(template.New("")).Parse(value))
}

// RegisterFuncs registers two new template functions with t.
// "JS" and "CSS" that return the appropriate resources to insert into the template.
func (resources *Resources) RegisterFuncs(t *template.Template) *template.Template {
	var builder strings.Builder
	resources.WriteCSS(&builder)
	css := template.HTML(builder.String())

	builder.Reset()
	resources.WriteJS(&builder)
	js := template.HTML(builder.String())

	return t.Funcs(template.FuncMap{
		"JS":  func() template.HTML { return js },
		"CSS": func() template.HTML { return css },
	})
}
