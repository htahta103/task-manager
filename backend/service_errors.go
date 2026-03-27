package main

import (
	"errors"
	"fmt"
)

var ErrInvalidUUID = errors.New("invalid uuid")

type InvalidIDError struct {
	Param   string
	Value   string
	Problem string // "missing" | "invalid"
}

func (e InvalidIDError) Error() string {
	param := e.Param
	if param == "" {
		param = "id"
	}
	switch e.Problem {
	case "missing":
		return fmt.Sprintf("missing %s", param)
	default:
		return fmt.Sprintf("invalid %s %q: expected UUID", param, e.Value)
	}
}

type validationError struct {
	msg string
}

func (e validationError) Error() string { return e.msg }

func ErrValidation(msg string) error { return validationError{msg: msg} }

func IsValidation(err error) bool {
	_, ok := err.(validationError)
	return ok
}
