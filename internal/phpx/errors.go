//spellchecker:words phpx
package phpx

import "fmt"

// Common PHP Errors.
const (
	errInit    = "Server initialization failed"
	errClosed  = "Server closed"
	errSend    = "Failed to encode request"
	errReceive = "Failed to decode response"
)

// PHPError represents an error during PHPServer logic.
type ServerError struct {
	Message string
	Err     error
}

// Unwrap returns the underlying error.
func (err ServerError) Unwrap() error {
	return err.Err
}

func (err ServerError) Error() string {
	if err.Err == nil {
		return "PHPServer: " + err.Message
	}
	return fmt.Sprintf("PHPServer: %s: %s", err.Message, err.Err)
}

// ThrowableError represents an error during php code.
type ThrowableError string

func (throwable ThrowableError) Error() string {
	return string(throwable)
}
