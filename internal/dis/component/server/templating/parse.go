//spellchecker:words templating
package templating

//spellchecker:words html template reflect golang slices
import (
	"html/template"
	"reflect"

	"golang.org/x/exp/slices"
)

// Parsed represents a parsed template that takes as argument a context of type C.
type Parsed[C any] struct {
	// does the context type an embed of the runtime flags type?
	hasRuntimeFlagsEmbed bool

	tpl   *template.Template // parsed template
	funcs []FlagFunc         // optionally concfigured functions.
}

// Parse parses a template with the given name and source.
// If base is not nil, every template associated with the base template is copied into the given template.
// Functions will be applied on creation time to represent the context for the given template.
func Parse[C any](name string, source []byte, base *template.Template, funcs ...FlagFunc) Parsed[C] {
	tp := reflect.TypeFor[C]()

	// determine if we have an embedded field in the struct
	var hasEmbed bool
	if tp.Kind() == reflect.Struct {
		field, ok := tp.FieldByName(runtimeFlagsName)
		if ok {
			hasEmbed = field.Anonymous
		}
	}

	// create a new template, and optionally inherit from the base template
	new := template.New(name)
	if base != nil {
		for _, tree := range base.Templates() {
			root := tree.Tree.Copy()
			if _, err := new.AddParseTree(tree.Name(), root); err != nil {
				panic("never reached") // Tree is a copy and has never been executed
			}
		}
	}

	return Parsed[C]{
		hasRuntimeFlagsEmbed: hasEmbed,
		tpl:                  template.Must(new.Parse(string(source))),
		funcs:                funcs,
	}
}

// Prepare prepares this template to be used with the given templating.
func (p *Parsed[C]) Prepare(templating *Templating, funcs ...FlagFunc) *Template[C] {
	pcopy := *p // make a copy of p!

	wrap := Template[C]{
		templating: templating,

		p: &pcopy,
	}

	// copy the functions!
	pcopy.funcs = slices.Clone(pcopy.funcs)
	pcopy.funcs = append(wrap.p.funcs, funcs...)

	return &wrap
}
