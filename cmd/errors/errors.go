// Package errors provides custom error types for the h3s project
package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"
)

// ErrorType is a type of error, specific to the h3s project
type ErrorType string

const (
	ErrorTypeConfig     ErrorType = "ConfigError"     // ErrorTypeConfig is an error type for configuration errors e.g. invalid configuration file
	ErrorTypeCluster    ErrorType = "ClusterError"    // ErrorTypeCluster is an error type for cluster errors
	ErrorTypeKubectl    ErrorType = "KubectlError"    // ErrorTypeKubectl is an error type for kubectl errors
	ErrorTypeHetzner    ErrorType = "HetznerError"    // ErrorTypeHetzner is an error type for Hetzner Cloud errors
	ErrorTypeK3s        ErrorType = "K3sError"        // ErrorTypeK3s is an error type for k3s errors
	ErrorTypeSSH        ErrorType = "SSHError"        // ErrorTypeSSH is an error type for ssh errors
	ErrorTypeSystem     ErrorType = "SystemError"     // ErrorTypeSystem is an error type for system errors
	ErrorTypeInternal   ErrorType = "InternalError"   // ErrorTypeInternal is an error type for internal errors
	ErrorTypeValidation ErrorType = "ValidationError" // ErrorTypeValidation is an error type for input validation errors
	ErrorTypeNetwork    ErrorType = "NetworkError"    // ErrorTypeNetwork is an error type for network-related errors
	ErrorTypeAuth       ErrorType = "AuthError"       // ErrorTypeAuth is an error type for authentication-related errors
	ErrorTypeResource   ErrorType = "ResourceError"   // ErrorTypeResource is an error type for resource management errors
	ErrorTypeTimeout    ErrorType = "TimeoutError"    // ErrorTypeTimeout is an error type for timeout-related errors
)

// Severity represents the severity level of an error
type Severity string

const (
	SeverityFatal   Severity = "FATAL"   // SeverityFatal indicates an unrecoverable error
	SeverityError   Severity = "ERROR"   // SeverityError indicates a serious but potentially recoverable error
	SeverityWarning Severity = "WARNING" // SeverityWarning indicates a warning that might need attention
)

// StackTrace represents a stack trace
type StackTrace struct {
	Function   string
	File       string
	SourceLine string
	Line       int
}

// Error is a custom error type for the h3s project
type Error struct {
	Timestamp time.Time
	Err       error
	Context   map[string]interface{}
	Message   string
	Operation string
	Type      ErrorType
	Severity  Severity
	Stack     []StackTrace
	Retryable bool
	Code      string // Error code for more specific error identification
}

// Error returns the string representation of the error
func (e *Error) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "[%s][%s] %s", e.Type, e.Severity, e.Message)
	if e.Operation != "" {
		fmt.Fprintf(&b, " (operation: %s)", e.Operation)
	}
	if e.Err != nil {
		fmt.Fprintf(&b, ": %v", e.Err)
	}
	if len(e.Context) > 0 {
		fmt.Fprintf(&b, "\nContext: %v", e.Context)
	}
	if len(e.Stack) > 0 {
		b.WriteString("\nStack Trace:")
		for _, frame := range e.Stack {
			fmt.Fprintf(&b, "\n\t%s:%d %s", frame.File, frame.Line, frame.Function)
		}
	}
	if e.Retryable {
		b.WriteString("\nThis operation may be retried")
	}
	return b.String()
}

// getStack returns the current stack trace
func getStack(skip int) []StackTrace {
	var stack []StackTrace
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			break
		}
		stack = append(stack, StackTrace{
			File:     file,
			Line:     line,
			Function: fn.Name(),
		})
		if len(stack) >= 10 { // Limit stack trace depth
			break
		}
	}
	return stack
}

// New creates a new error with the given type and message
func New(errType ErrorType, message string) *Error {
	return &Error{
		Type:      errType,
		Message:   message,
		Stack:     getStack(2),
		Context:   make(map[string]interface{}),
		Severity:  SeverityError,
		Retryable: false,
		Timestamp: time.Now(),
	}
}

// Wrap wraps an error with a custom error type and message
func Wrap(errType ErrorType, message string, err error) *Error {
	return &Error{
		Type:      errType,
		Message:   message,
		Err:       err,
		Stack:     getStack(2),
		Context:   make(map[string]interface{}),
		Severity:  SeverityError,
		Retryable: false,
		Timestamp: time.Now(),
	}
}

// WithCode adds an error code to an error
func (e *Error) WithCode(code string) *Error {
	e.Code = code
	return e
}

// IsErrorType checks if an error is of a specific ErrorType
func IsErrorType(err error, errType ErrorType) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Type == errType
	}
	return false
}

// GetErrorContext extracts context from an error if it's a custom Error
func GetErrorContext(err error) map[string]interface{} {
	var e *Error
	if errors.As(err, &e) {
		return e.Context
	}
	return nil
}

// WithContext adds context to an error
func (e *Error) WithContext(key string, value interface{}) *Error {
	e.Context[key] = value
	return e
}

// WithOperation adds operation information to an error
func (e *Error) WithOperation(operation string) *Error {
	e.Operation = operation
	return e
}

// WithSeverity sets the severity level of the error
func (e *Error) WithSeverity(severity Severity) *Error {
	e.Severity = severity
	return e
}

// WithRetryable marks whether the error is retryable
func (e *Error) WithRetryable(retryable bool) *Error {
	e.Retryable = retryable
	return e
}

// Is reports whether err is of type target
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.Type == t.Type
}

// Unwrap returns the wrapped error
func (e *Error) Unwrap() error {
	return e.Err
}
