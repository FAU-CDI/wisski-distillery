package phpx

import "fmt"

// Common PHP Errors
const (
	errInit    = "Server initialization failed"
	errClosed  = "Server closed"
	errReceive = "Failed to decode response"
)

// PHPError represents an error during PHPServer logic
type ServerError struct {
	Message string
	Err     error
}

// Unwrap returns the underlying error
func (err ServerError) Unwrap() error {
	return err.Err
}

func (err ServerError) Error() string {
	if err.Err == nil {
		return fmt.Sprintf("PHPServer: %s", err.Message)
	}
	return fmt.Sprintf("PHPServer: %s: %s", err.Message, err.Err)
}

// Throwable represents an error during php code
type Throwable string

func (throwable Throwable) Error() string {
	return string(throwable)
}
