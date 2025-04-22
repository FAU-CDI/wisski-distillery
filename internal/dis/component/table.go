//spellchecker:words component
package component

//spellchecker:words reflect
import (
	"reflect"
)

// Table is a component that manages a provided sql table.
type Table interface {
	Component

	// TableInfo returns information about a specific table
	TableInfo() TableInfo
}

type TableInfo struct {
	Model reflect.Type // model is the model this type manages
	Name  string
}
