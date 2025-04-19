// helpers.go - Helper functions for configuring mock behaviors in tests
package mockhetzner

import (
	"github.com/stretchr/testify/mock"
)

// SetMockReturn sets up the mock to return the given values for a method.
func SetMockReturn(m *mock.Mock, method string, returns ...interface{}) {
	m.On(method, mock.Anything).Return(returns...)
}

// SetMockReturnWithArgs sets up the mock to return values for a method with specific arguments.
func SetMockReturnWithArgs(m *mock.Mock, method string, args []interface{}, returns ...interface{}) {
	m.On(method, args...).Return(returns...)
}
