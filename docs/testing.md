# H3S Testing Guide

This document provides detailed information about the testing approach used in H3S, including test organization, test data management, mocking strategies, and best practices.

## Test Organization

H3S tests are organized into the following categories:

### Unit Tests

Unit tests focus on testing individual functions and methods in isolation. They are located alongside the code they test with the `_test.go` suffix.

```
internal/
  config/
    config.go
    config_test.go  # Unit tests for config package
  hetzner/
    hetzner.go
    hetzner_test.go  # Unit tests for hetzner package
```

### Integration Tests

Integration tests verify that different components work together correctly. They are located in the `/test/integration` directory.

```
test/
  integration/
    config_validation_test.go  # Tests config validation with real file I/O
    core_commands_test.go      # Tests CLI commands
```

### End-to-End Tests

End-to-end tests verify complete workflows from the user's perspective. They are also located in the `/test/integration` directory but focus on testing entire workflows.

```
test/
  integration/
    e2e_test.go                # Tests complete workflows
    workflow_modification_test.go  # Tests config modification workflows
```

### Benchmarks

Benchmarks measure the performance of critical operations. They are located alongside the code they benchmark with the `_test.go` suffix and function names starting with `Benchmark`.

```
internal/
  config/
    benchmark_test.go  # Benchmarks for config operations
  hetzner/
    resource_benchmark_test.go  # Benchmarks for resource operations
```

## Running Tests

### Running All Tests

```bash
make test
```

### Running Unit Tests Only

```bash
make unit-test
```

### Running Integration Tests Only

```bash
make integration-test
```

### Running End-to-End Tests Only

```bash
make e2e-test
```

### Running a Specific Test

```bash
make unit-test NAME=TestConfigValidation
make integration-test NAME=TestConfigValidationIntegration
make e2e-test NAME=TestE2EWorkflows
```

### Running Benchmarks

```bash
make bench
```

### Running a Specific Benchmark

```bash
make bench NAME=BenchmarkConfigLoad
```

## Test Data Management

H3S uses several approaches for managing test data:

### Temporary Files and Directories

For tests that need to work with files, we use the `os.MkdirTemp()` and `os.CreateTemp()` functions to create temporary files and directories that are automatically cleaned up after the test.

```go
func TestSaveConfig(t *testing.T) {
    // Create a temporary directory for the test
    tempDir, err := os.MkdirTemp("", "h3s-test-*")
    if err != nil {
        t.Fatalf("Failed to create temp directory: %v", err)
    }
    defer os.RemoveAll(tempDir)
    
    // Use the temporary directory for test files
    configPath := filepath.Join(tempDir, "h3s.yaml")
    // ... test implementation
}
```

### Test Fixtures

For complex test data, we use fixtures in the `testdata` directory. These are static files that are used by tests.

```go
func TestLoadConfig(t *testing.T) {
    // Load a test fixture
    configPath := "testdata/valid_config.yaml"
    config, err := config.Load(configPath)
    // ... test implementation
}
```

### Helper Functions

We use helper functions to generate test data programmatically. These functions create consistent test data for use in multiple tests.

```go
// createTestConfig creates a test configuration file
func createTestConfig(dir string, mockServerURL string) (string, error) {
    configPath := filepath.Join(dir, "h3s.yaml")
    config := []byte(`
name: test-cluster
k3s_version: v1.28.2+k3s1
# ... other required fields
hetzner_api_endpoint: "` + mockServerURL + `"
`)
    return configPath, os.WriteFile(configPath, config, 0644)
}
```

## Mocking Approach

H3S uses several approaches for mocking external dependencies:

### Mock Hetzner API

The `mockhetzner` package provides a mock implementation of the Hetzner Cloud API. It simulates API responses for testing without making real API calls.

```go
// Create a mock server
mockServer := mockhetzner.NewHetznerMockScenario("/v1", "success")
defer mockServer.Close()

// Use the mock server URL in tests
config.HetznerAPIEndpoint = mockServer.Server.URL
```

### Interface-Based Mocking

We define interfaces for external dependencies and create mock implementations for testing.

```go
// Define an interface
type CloudProvider interface {
    CreateServer(name string, options ServerOptions) (*Server, error)
    DeleteServer(id string) error
}

// Create a mock implementation
type MockCloudProvider struct {
    CreateServerFunc func(name string, options ServerOptions) (*Server, error)
    DeleteServerFunc func(id string) error
}

func (m *MockCloudProvider) CreateServer(name string, options ServerOptions) (*Server, error) {
    return m.CreateServerFunc(name, options)
}

func (m *MockCloudProvider) DeleteServer(id string) error {
    return m.DeleteServerFunc(id)
}
```

### Dependency Injection

We use dependency injection to pass mock implementations to functions and methods.

```go
// Function that accepts a dependency
func CreateCluster(ctx *cluster.Cluster, provider CloudProvider) error {
    // Use the provider
    server, err := provider.CreateServer("control-plane", options)
    // ... implementation
}

// Test with a mock
func TestCreateCluster(t *testing.T) {
    mockProvider := &MockCloudProvider{
        CreateServerFunc: func(name string, options ServerOptions) (*Server, error) {
            return &Server{ID: "123", Name: name}, nil
        },
    }
    
    ctx := &cluster.Cluster{/* ... */}
    err := CreateCluster(ctx, mockProvider)
    // ... assertions
}
```

### Function Variable Replacement

For functions that are difficult to mock with interfaces, we use function variables that can be replaced in tests.

```go
// Define a function variable
var execCommand = exec.Command

// Function that uses the variable
func RunCommand(name string, args ...string) (string, error) {
    cmd := execCommand(name, args...)
    output, err := cmd.CombinedOutput()
    return string(output), err
}

// Test with a mock
func TestRunCommand(t *testing.T) {
    // Save the original
    origExecCommand := execCommand
    defer func() { execCommand = origExecCommand }()
    
    // Replace with a mock
    execCommand = func(name string, args ...string) *exec.Cmd {
        return exec.Command("echo", "mock output")
    }
    
    output, err := RunCommand("ls", "-la")
    // ... assertions
}
```

## Best Practices

### Writing Good Tests

1. **Test one thing at a time**: Each test should focus on a single aspect of functionality.
2. **Use descriptive test names**: Test names should describe what is being tested and the expected outcome.
3. **Set up and tear down properly**: Create any necessary test data before the test and clean up afterward.
4. **Use subtests for related tests**: Use `t.Run()` to group related tests together.
5. **Test edge cases and error conditions**: Don't just test the happy path.

### Test Coverage

We aim for high test coverage to ensure code quality and reliability:

- **Core business logic**: 80%+ coverage
- **CLI commands**: 70%+ coverage
- **Utility functions**: 60%+ coverage

To check coverage locally:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Continuous Integration

All tests are run in CI for every pull request. The CI pipeline includes:

1. Running unit tests
2. Running integration tests
3. Running end-to-end tests
4. Measuring test coverage
5. Reporting coverage to Codecov

### Test Environment Variables

The following environment variables can be used to control test behavior:

- `H3S_ENABLE_E2E_TESTS=1`: Enable end-to-end tests (disabled by default)
- `H3S_ENABLE_REAL_INTEGRATION=1`: Enable tests that interact with real Hetzner APIs (disabled by default)
- `H3S_NON_INTERACTIVE=1`: Disable interactive prompts in tests

## Adding New Tests

When adding new features or fixing bugs, follow these steps to add tests:

1. **Add unit tests** for new functions and methods
2. **Add integration tests** for interactions between components
3. **Add end-to-end tests** for new workflows
4. **Add benchmarks** for performance-critical code

Example of adding a new unit test:

```go
func TestNewFeature(t *testing.T) {
    // Setup
    input := "test input"
    
    // Execute
    result, err := NewFeature(input)
    
    // Assert
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }
    if result != "expected output" {
        t.Errorf("Expected 'expected output', got: %q", result)
    }
}
```

Example of adding a new integration test:

```go
//go:build integration
// +build integration

func TestNewFeatureIntegration(t *testing.T) {
    // Setup mock server
    mockServer := mockhetzner.NewHetznerMockScenario("/v1", "success")
    defer mockServer.Close()
    
    // Create test config
    configPath, cleanup := createTestConfig(t, mockServer.Server.URL)
    defer cleanup()
    
    // Execute the command
    out, err := runCLI("new-feature", "--config", configPath)
    
    // Assert
    if err != nil {
        t.Errorf("Expected no error, got: %v\nOutput: %s", err, out)
    }
    if !strings.Contains(out, "expected output") {
        t.Errorf("Expected output to contain 'expected output', got: %q", out)
    }
}
```

## Troubleshooting Tests

### Common Issues

1. **Tests fail intermittently**: This may indicate a race condition or dependency on external state. Make sure tests are independent and deterministic.
2. **Tests are slow**: Consider using mocks instead of real external services. Use the `-short` flag to skip slow tests during development.
3. **Tests fail in CI but pass locally**: Check for environment differences, file paths, or timing issues.

### Debugging Tests

1. **Use `t.Log()` and `t.Logf()`**: Add logging to see what's happening during the test.
2. **Run tests with verbose output**: Use `go test -v` to see detailed output.
3. **Run a specific test**: Use `go test -run TestName` to run only the failing test.
4. **Use the debugger**: Use Delve or your IDE's debugger to step through the test.
