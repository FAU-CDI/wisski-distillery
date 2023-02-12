package validators

import "time"

func ValidateDuration(d *time.Duration, dflt string) error {
	if *d == 0 {
		var err error
		*d, err = time.ParseDuration(dflt)
		if err != nil {
			return err
		}
	}
	return nil
}
