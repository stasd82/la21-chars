package validate

import (
	"encoding/json"
	"errors"
)

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type FieldErrors []FieldError

func (f FieldErrors) Error() string {
	d, err := json.Marshal(f)
	if err != nil {
		return err.Error()
	}
	return string(d)
}

func IsFieldErrors(err error) bool {
	var fe FieldErrors
	return errors.As(err, &fe)
}

func GetFieldErrors(err error) FieldErrors {
	var fe FieldErrors
	if !errors.As(err, &fe) {
		return nil
	}
	return fe
}

func (f FieldErrors) Fields() map[string]string {
	m := make(map[string]string)
	for _, fld := range f {
		m[fld.Field] = fld.Error
	}
	return m
}
