package api

import (
	"fmt"
	"net/http"
)

// ErrorMessage is the message of an error response.
type ErrorMessage struct {
	Message string `json:"message"`
}

// UnauthorizedError represents the message of an HTTP 401 response.
type UnauthorizedError ErrorMessage

// UnprocessableEntityError represents the generic structure of an error response.
type UnprocessableEntityError struct {
	Error ErrorMessage `json:"error"`
}

// parseUnauthorizedError parses an HTTP 401 response into an UnauthorizedError.
func parseUnauthorizedError(resp *http.Response) (*UnauthorizedError, error) {
	var unauthorizedError UnauthorizedError
	if err := readAndParseJSONBody(resp, &unauthorizedError); err != nil {
		return nil, fmt.Errorf("error parsing unauthorized error response: %w", err)
	}
	return &unauthorizedError, nil
}

// parseUnprocessableEntityError parses an HTTP 422 response into an UnprocessableEntityError.
func parseUnprocessableEntityError(resp *http.Response) (*UnprocessableEntityError, error) {
	var unprocessableEntityError UnprocessableEntityError
	if err := readAndParseJSONBody(resp, &unprocessableEntityError); err != nil {
		return nil, fmt.Errorf("error parsing unprocessable entity error response: %w", err)
	}
	return &unprocessableEntityError, nil
}
