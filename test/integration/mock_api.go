package integration

import (
	"net/http"
	"net/http/httptest"
)

// MockAPIServer represents a simple mock API server for integration tests.
type MockAPIServer struct {
	Server *httptest.Server
}

// NewMockAPIServer creates and starts a new mock API server with the provided handler.
func NewMockAPIServer(handler http.Handler) *MockAPIServer {
	ts := httptest.NewServer(handler)
	return &MockAPIServer{Server: ts}
}

// Close shuts down the mock API server.
func (m *MockAPIServer) Close() {
	m.Server.Close()
}

// URL returns the base URL of the mock API server.
func (m *MockAPIServer) URL() string {
	return m.Server.URL
}
