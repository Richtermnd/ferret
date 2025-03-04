package object

import (
	"fmt"
)

type ErrorType string

const (
	UNSUPPORTED_ERR      = "unsupported"
	UNKNOWN_OPERATOR_ERR = "unknown operator"
	NOT_IMPLEMENTED_ERR  = "not implemented"
	NOT_FOUND_ERR        = "not found"
	UNEXPECTED           = "unexpected"
)

type Error struct {
	ErrType ErrorType
	msg     string
}

func (err *Error) Type() ObjectType { return ERROR_OBJ }
func (err *Error) Inspect() string  { return "[ERROR] " + string(err.ErrType) + ": " + err.msg }

func IsError(obj Object) bool {
	return obj.Type() == ERROR_OBJ
}

func NewError(t ErrorType, format string, args ...any) Object {
	return &Error{ErrType: t, msg: fmt.Sprintf(format, args...)}
}
