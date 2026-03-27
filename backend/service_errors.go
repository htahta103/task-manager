package main

import "errors"

var ErrInvalidUUID = errors.New("invalid uuid")

type validationError struct {
	msg string
}

func (e validationError) Error() string { return e.msg }

func ErrValidation(msg string) error { return validationError{msg: msg} }

func IsValidation(err error) bool {
	_, ok := err.(validationError)
	return ok
}
