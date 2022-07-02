package tux

import "errors"

type shutdownError struct {
	Message string
}

func NewShutdownError(msg string) error {
	return &shutdownError{Message: msg}
}

func (s *shutdownError) Error() string {
	return s.Message
}

func IsShutdown(err error) bool {
	var se *shutdownError
	return errors.As(err, &se)
}
