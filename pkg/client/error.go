package client

import (
	"fmt"
	"net/http"
)

type ErrHttpStatus struct {
	StatusCode int
}

func (err *ErrHttpStatus) Error() string {
	return fmt.Sprintf("http status %d", err.StatusCode)
}

func (err *ErrHttpStatus) IsAuthStatus() bool {
	return IsAuthErrorStatusCode(err.StatusCode)
}

func IsAuthErrorStatusCode(statusCode int) bool {
	return statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden
}

type ErrHttpBodyUnmarshal struct {
	Cause error
}

func (err *ErrHttpBodyUnmarshal) Error() string {
	return fmt.Sprintf("http body unmarshal error: %s", err.Cause.Error())
}

func (err *ErrHttpBodyUnmarshal) Unwrap() error {
	return err.Cause
}
