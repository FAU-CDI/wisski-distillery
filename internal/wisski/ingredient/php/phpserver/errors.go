package phpserver

import "fmt"

// Common PHP Errors
var (
	errPHPInit    = "Unable to initialize"
	errPHPMarshal = "Marshal failed"
	errPHPInvalid = ServerError{Message: "Invalid code to execute"}
	errPHPReceive = "Failed to receive response"
	errPHPClosed  = ServerError{Message: "Server closed"}
)

// PHPError represents an error during PHPServer logic
type ServerError struct {
	Message string
	Err     error
}

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
