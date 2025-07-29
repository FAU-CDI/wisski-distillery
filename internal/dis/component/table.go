//spellchecker:words component
package component

//spellchecker:words github wisski distillery internal models
import "github.com/FAU-CDI/wisski-distillery/internal/models"

//spellchecker:words reflect

// Table is a component that manages a provided sql table.
type Table interface {
	Component

	// TableInfo returns information about a specific table
	TableInfo() TableInfo
}

type TableInfo struct {
	Model models.Model
}

func (ti TableInfo) Name() string {
	return ti.Model.TableName()
}
