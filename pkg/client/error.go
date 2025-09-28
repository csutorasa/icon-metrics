package client

import (
	"fmt"
)

// Error that is used when the iCon device responds with an unsuccessful status code.
type ErrHttpStatus struct {
	// HTTP response status code.
	StatusCode int
}

// implements error
func (err *ErrHttpStatus) Error() string {
	return fmt.Sprintf("http status %d", err.StatusCode)
}

// Error that is used when the iCon device responds with a HTTP body that the client cannot parse.
type ErrHttpBodyUnmarshal struct {
	Cause error
}

// implements error
func (err *ErrHttpBodyUnmarshal) Error() string {
	return fmt.Sprintf("http body unmarshal error: %s", err.Cause.Error())
}

// Unwraps the cause.
func (err *ErrHttpBodyUnmarshal) Unwrap() error {
	return err.Cause
}
