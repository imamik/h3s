// Package mockhetzner provides mock implementations for Hetzner API interfaces.
package mockhetzner

import (
	"github.com/stretchr/testify/mock"
)

// API is a mock of the Hetzner API client interface.
type API struct {
	mock.Mock
}

// Example mocked method for demonstration
func (m *API) CreateServer(name string) (string, error) {
	args := m.Called(name)
	return args.String(0), args.Error(1)
}
