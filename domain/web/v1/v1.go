package v1

import "errors"

type ErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

type RequestErr struct {
	Err    error
	Status int
}

func NewRequestErr(err error, status int) error {
	return &RequestErr{Err: err, Status: status}
}

func (r *RequestErr) Error() string {
	return r.Err.Error()
}

func IsRequestErr(err error) bool {
	var re *RequestErr
	return errors.As(err, &re)
}

func GetRequestErr(err error) *RequestErr {
	var re *RequestErr
	if !errors.As(err, &re) {
		return nil
	}
	return re
}
