//spellchecker:words component
package component

//spellchecker:words context
import (
	"context"
)

// Cronable is a component that implements a cron method.
type Cronable interface {
	Component

	// Name of the cron task being performed
	TaskName() string

	// Cron is called to run this cron task
	Cron(ctx context.Context) error
}
