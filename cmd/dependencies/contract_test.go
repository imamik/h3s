//go:build test
// +build test

package dependencies

import "testing"

// contractImplementations lists all types that must satisfy CommandDependencies.
// Add new implementations here to enforce the contract.
var contractImplementations = []CommandDependencies{
	// Example: &MyDependencies{},
}

// TestCommandDependenciesContract ensures all listed implementations satisfy the CommandDependencies interface at compile time.
func TestCommandDependenciesContract(t *testing.T) {
	for _, impl := range contractImplementations {
		_ = impl // Compile-time check: will not compile if type does not implement interface
	}
}

// To add runtime checks for method presence, use reflection as needed.
// For more advanced contract enforcement, consider static analysis or code generation tools.
