//spellchecker:words validators
package validators

//spellchecker:words time
import (
	"fmt"
	"time"
)

func ValidateDuration(d *time.Duration, dflt string) error {
	if *d == 0 {
		var err error
		*d, err = time.ParseDuration(dflt)
		if err != nil {
			return fmt.Errorf("failed to parse duration: %w", err)
		}
	}
	return nil
}
