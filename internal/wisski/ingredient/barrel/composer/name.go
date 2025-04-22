//spellchecker:words composer
package composer

//spellchecker:words strings
import "strings"

// If the module name cannot be found, returns the string unchanged.
func ModuleName(spec string) string {
	_, name, found := strings.Cut(spec, "/")
	if !found {
		return spec
	}
	name, _, _ = strings.Cut(name, ":")
	return name
}
