// Package errors provides custom error types for the h3s project
package errors

import (
	"fmt"
	"strings"
)

// ErrorType is a type of error, specific to the h3s project
type ErrorType string

const (
	ErrorTypeConfig  ErrorType = "ConfigError"  // ErrorTypeConfig is an error type for configuration errors e.g. invalid configuration file
	ErrorTypeCluster ErrorType = "ClusterError" // ErrorTypeCluster is an error type for cluster errors
	ErrorTypeKubectl ErrorType = "KubectlError" // ErrorTypeKubectl is an error type for kubectl errors
	ErrorTypeHetzner ErrorType = "HetznerError" // ErrorTypeHetzner is an error type for Hetzner Cloud errors
	ErrorTypeK3s     ErrorType = "K3sError"     // ErrorTypeK3s is an error type for k3s errors
	ErrorTypeSSH     ErrorType = "SSHError"     // ErrorTypeSSH is an error type for ssh errors
	ErrorTypeSystem  ErrorType = "SystemError"  // ErrorTypeSystem is an error type for system errors
)

// Error is a custom error type for the h3s project
type Error struct {
	Message string
	Err     error
	Type    ErrorType
}

// Error returns the string representation of the error
func (e *Error) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "[%s] %s", e.Type, e.Message)
	if e.Err != nil {
		fmt.Fprintf(&b, ": %v", e.Err)
	}
	return b.String()
}

// New creates a new error with the given type and message
func New(errType ErrorType, message string) *Error {
	return &Error{
		Type:    errType,
		Message: message,
	}
}

// Wrap wraps an error with a custom error type and message
func Wrap(errType ErrorType, message string, err error) *Error {
	return &Error{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}
