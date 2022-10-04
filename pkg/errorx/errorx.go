package errorx

// First returns the first non-nil error, or nil otherwise.
func First(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}
