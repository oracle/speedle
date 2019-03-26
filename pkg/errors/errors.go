//Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package errors

import (
	"fmt"

	gerrs "github.com/pkg/errors"
)

// speedleError is the struct of speedle error
type speedleError struct {
	code    ErrorCode
	message string
	cause   error
}

func (e *speedleError) Error() string {
	errMsg := fmt.Sprintf("%s %s", e.code, e.message)
	if nil == e.cause {
		return errMsg
	}

	return errMsg + ": " + e.cause.Error()
}

func (e *speedleError) Cause() error {
	return e.cause
}

func (e *speedleError) Code() ErrorCode {
	return e.code
}

// Cause returns the cause error of this error
func Cause(err error) error {
	return gerrs.Cause(err)
}

// Code returns the error code
func Code(err error) ErrorCode {
	type coder interface {
		Code() ErrorCode
	}

	cd, ok := err.(coder)
	if !ok {
		return UnknownError
	}
	return cd.Code()
}

// Errorf formats an error with format
func Errorf(code ErrorCode, format string, a ...interface{}) error {
	return &speedleError{
		code:    code,
		message: fmt.Sprintf(format, a...),
	}
}

// New constructs a new error
func New(code ErrorCode, message string) error {
	return &speedleError{
		code:    code,
		message: message,
	}
}

// Wrapf warps an error with a error code and a format message
func Wrapf(err error, code ErrorCode, format string, a ...interface{}) error {
	return Wrap(err, code, fmt.Sprintf(format, a...))
}

// Wrap waps an error with an error and a message
func Wrap(err error, code ErrorCode, message string) error {
	if err == nil {
		return nil
	}
	return &speedleError{
		code:    code,
		message: message,
		cause:   err,
	}
}
