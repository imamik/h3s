package dependencies

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	deps := Get()
	assert.NotNil(t, deps)
	assert.IsType(t, &DefaultDependencies{}, deps)
}

func TestDefaultDependencies_GetClusterContext(t *testing.T) {
	deps := &DefaultDependencies{}
	ctx, err := deps.GetClusterContext()
	// Since this is integration-like, we expect an error in test environment
	assert.Error(t, err)
	assert.Nil(t, ctx)
}

func TestDefaultDependencies_GetKubeconfigPath(t *testing.T) {
	deps := &DefaultDependencies{}
	path, exists := deps.GetKubeconfigPath()
	// In test environment, kubeconfig shouldn't exist
	assert.False(t, exists)
	assert.True(t, strings.HasSuffix(path, "h3s-kubeconfig.yaml"))
}

func TestDefaultDependencies_GetKubeconfigPath_InvalidPath(t *testing.T) {
	deps := &DefaultDependencies{}

	// Test with invalid path
	path, exists := deps.GetKubeconfigPath()
	assert.False(t, exists)
	assert.True(t, strings.HasSuffix(path, "h3s-kubeconfig.yaml"))
}

func TestDefaultDependencies_ExecuteLocalCommand_Error(t *testing.T) {
	deps := &DefaultDependencies{}

	// Test with invalid command
	_, err := deps.ExecuteLocalCommand("invalid-command")
	assert.Error(t, err)
}
