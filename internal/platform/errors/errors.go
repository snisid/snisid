package errors

import (
	"bytes"
	"fmt"
	"runtime/debug"
)

// Error is the standard SNISID platform error.
type Error struct {
	Code    ErrorCode
	Message string
	Op      string
	Err     error
	Trace   []byte
}

func (e *Error) Error() string {
	var buf bytes.Buffer

	if e.Op != "" {
		fmt.Fprintf(&buf, "[%s] ", e.Op)
	}

	if e.Code != "" {
		fmt.Fprintf(&buf, "<%s> ", e.Code)
	}

	if e.Message != "" {
		buf.WriteString(e.Message)
	}

	if e.Err != nil {
		buf.WriteString(": ")
		buf.WriteString(e.Err.Error())
	}

	return buf.String()
}

// Unwrap allows standard errors.Is and errors.As usage.
func (e *Error) Unwrap() error {
	return e.Err
}

// New creates a new structured Error. Stack traces are attached automatically for Internal errors.
func New(code ErrorCode, message string, op string, err error) *Error {
	e := &Error{
		Code:    code,
		Message: message,
		Op:      op,
		Err:     err,
	}

	if code == Internal {
		e.Trace = debug.Stack()
	}

	return e
}
