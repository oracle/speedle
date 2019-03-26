//Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package errors

// ErrorCode is data type of error codes for different kind of errors
type ErrorCode string

// UnknownError is the unknown error
const UnknownError ErrorCode = "SPDL-0000"

// For common components
const (
	ConfigError    ErrorCode = "SPDL-0001"
	ServerError    ErrorCode = "SPDL-0002"
	LoggingError   ErrorCode = "SPDL-0003"
	InvalidRequest ErrorCode = "SPDL-0004"
)

// For policy management errors
const (
	StoreError          ErrorCode = "SPDL-1001"
	EntityNotFound      ErrorCode = "SPDL-1002"
	EntityAlreadyExists ErrorCode = "SPDL-1003"
	ExceedLimit         ErrorCode = "SPDL-1004"
	SerializationError  ErrorCode = "SPDL-1005"
)

// For evaluator errors
const (
	EvalEngineError   ErrorCode = "SPDL-2001"
	EvalCacheError    ErrorCode = "SPDL-2002"
	BuiltInFuncError  ErrorCode = "SPDL-2003"
	CustomerFuncError ErrorCode = "SPDL-2004"
	DiscoverError     ErrorCode = "SPDL-2005"
)
