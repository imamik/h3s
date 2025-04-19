// Package mockhetzner provides mock implementations for Hetzner API interfaces.
package mockhetzner

import (
	"github.com/stretchr/testify/mock"
)

// MockHetznerAPI is a mock of the Hetzner API client interface.
type MockHetznerAPI struct {
	mock.Mock
}

// Example mocked method for demonstration
func (m *MockHetznerAPI) CreateServer(name string) (string, error) {
	args := m.Called(name)
	return args.String(0), args.Error(1)
}
