// Package resources provides Resources
package resources

import (
	"html/template"
	"io"
	"strings"

	"github.com/tkw1536/goprogram/lib/collection"
	"golang.org/x/net/html"
)

// Resources represents resources found inside a "html" file
type Resources struct {
	JSModules []string // <script type="module">
	JSRegular []string // <script>
	CSS       []string // <link rel="stylesheet">
}

var attributeEscaper = strings.NewReplacer("<", "&lt;", "&", "&amp;", "\"", "&quot;")

const attributeShouldQuote = "<&\" "
const quoteString = "\""

func attributeValue(value string) string {
	value = attributeEscaper.Replace(value)
	if strings.ContainsAny(value, attributeShouldQuote) {
		return quoteString + value + quoteString
	}
	return value
}

var openLinkBytes = []byte("<link rel=stylesheet href=")
var closeLinkBytes = []byte(">")
var openModuleBytes = []byte("<script type=module src=")
var openRegularBytes = []byte("<script src=")
var closeScriptBytes = []byte("></script>")

// WriteCSS writes all link tags to writer
func (resources *Resources) WriteCSS(writer io.Writer) {
	for _, href := range resources.CSS {
		writer.Write(openLinkBytes)
		writer.Write([]byte(attributeValue(href)))
		writer.Write(closeLinkBytes)
	}
}

func (resources *Resources) CSSTemplate() template.HTML {
	var buffer strings.Builder
	resources.WriteCSS(&buffer)
	return template.HTML(buffer.String())
}

// WriteJS writes all JavaScript tags to writer
func (resources *Resources) WriteJS(writer io.Writer) {
	for _, href := range resources.JSModules {
		writer.Write(openModuleBytes)
		writer.Write([]byte(attributeValue(href)))
		writer.Write(closeScriptBytes)
	}
	for _, href := range resources.JSRegular {
		writer.Write(openRegularBytes)
		writer.Write([]byte(attributeValue(href)))
		writer.Write(closeScriptBytes)
	}
}

func (resources *Resources) JSTemplate() template.HTML {
	var buffer strings.Builder
	resources.WriteJS(&buffer)
	return template.HTML(buffer.String())
}

// Parse parses resources from reader
func Parse(r io.Reader) (src Resources) {
	z := html.NewTokenizer(r)
	for {
		// read the next token
		z.Next()
		token := z.Token()

		switch {
		case token.Type == html.ErrorToken:
			return

		case token.Type == html.StartTagToken && token.Data == "script":
			// <script src="...">
			js := getAttributeValue(token, "src")
			if js != "" {
				if getAttributeValue(token, "type") == "module" {
					src.JSModules = append(src.JSModules, js)
				} else {
					src.JSRegular = append(src.JSRegular, js)
				}
			}
		}
		if token.Type == html.StartTagToken && token.Data == "link" && getAttributeValue(token, "rel") == "stylesheet" {
			// <link rel="stylesheet" href="...">
			css := getAttributeValue(token, "href")
			if css != "" {
				src.CSS = append(src.CSS, css)
			}
		}
	}

}

// getAttributeValue returns the value of the given attribute, or the empty string if it is unset
func getAttributeValue(token html.Token, attr string) string {
	return collection.First(token.Attr, func(a html.Attribute) bool {
		return a.Key == attr
	}).Val
}
