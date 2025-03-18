package myerrs

import "fmt"

type StatusCodeError struct {
	StatusCode int
}

func (sc *StatusCodeError) Error() string {
	return fmt.Sprintf("status code: %d", sc.StatusCode)
}

type PlatformError struct {
	Text string
}

func (pe *PlatformError) Error() string {
	return fmt.Sprintf("platform returned error: (%s)", pe.Text)
}
