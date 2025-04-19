# Contributing to H3S

Thank you for your interest in contributing to H3S! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Testing Standards](#testing-standards)
  - [Test Organization](#test-organization)
  - [Test Types](#test-types)
  - [Writing Good Tests](#writing-good-tests)
  - [Test Data Management](#test-data-management)
  - [Mocking Approach](#mocking-approach)
  - [Coverage Expectations](#coverage-expectations)
- [Pull Request Process](#pull-request-process)
- [Style Guidelines](#style-guidelines)
- [Documentation](#documentation)

## Code of Conduct

Please be respectful and considerate of others when contributing to this project. We aim to foster an inclusive and welcoming community.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/h3s.git`
3. Set up the development environment: `make dev`
4. Create a new branch for your feature or bugfix: `git checkout -b feature/your-feature-name`

## Development Workflow

1. Make your changes
2. Run tests to ensure your changes don't break existing functionality: `make test`
3. Run linting to ensure code quality: `golangci-lint run`
4. Commit your changes following the [Conventional Commits](https://www.conventionalcommits.org/) style
5. Push your changes to your fork
6. Submit a pull request

## Testing Standards

H3S follows a comprehensive testing approach to ensure code quality and reliability. This section outlines the testing standards and expectations for contributors.

### Test Organization

Tests in H3S are organized into the following directories:

- **Unit tests**: Located alongside the code they test with `_test.go` suffix
- **Integration tests**: Located in `/test/integration`
- **End-to-end tests**: Located in `/test/integration` with specific test files

### Test Types

#### Unit Tests

Unit tests focus on testing individual functions and methods in isolation. They should:

- Be fast and not depend on external services
- Use mocks for external dependencies
- Cover edge cases and error conditions
- Be named with the pattern `Test<FunctionName>_<Scenario>`

Example:

```go
func TestValidateConfig_MissingRequiredField(t *testing.T) {
    // Test implementation
}
```

#### Integration Tests

Integration tests verify that different components work together correctly. They:

- May use mocked external services (like Hetzner API)
- Test interactions between multiple components
- Verify correct API usage patterns
- Should be tagged with `//go:build integration`

Example:

```go
//go:build integration
// +build integration

func TestConfigValidationIntegration(t *testing.T) {
    // Test implementation
}
```

#### End-to-End Tests

End-to-end tests verify complete workflows from the user's perspective. They:

- Test the CLI interface and commands
- Verify entire workflows (create, get, destroy)
- Use mocked external services by default
- Can optionally use real services when explicitly enabled
- Should be tagged with `//go:build e2e`

Example:

```go
//go:build e2e
// +build e2e

func TestE2EWorkflows(t *testing.T) {
    // Test implementation
}
```

#### Benchmarks

Benchmarks measure the performance of critical operations. They:

- Focus on performance-critical code paths
- Should be named with the pattern `Benchmark<Operation>`
- Can include sub-benchmarks for different scenarios
- Should use realistic data sizes

Example:

```go
func BenchmarkConfigLoad(b *testing.B) {
    // Benchmark implementation
}
```

### Writing Good Tests

Good tests in H3S should:

1. **Be clear and readable**: Test names should describe what is being tested and the expected outcome
2. **Be independent**: Tests should not depend on each other or on external state
3. **Be deterministic**: Tests should produce the same result every time they run
4. **Test one thing**: Each test should focus on a single aspect of functionality
5. **Include assertions**: Tests should verify expected outcomes, not just run code
6. **Handle cleanup**: Tests should clean up any resources they create

Example of a good test:

```go
func TestCreateConfig_ValidInput(t *testing.T) {
    // Setup
    tempDir, err := os.MkdirTemp("", "h3s-test-*")
    if err != nil {
        t.Fatalf("Failed to create temp directory: %v", err)
    }
    defer os.RemoveAll(tempDir)
    
    // Test execution
    config := &config.Config{
        Name: "test-cluster",
        // ... other required fields
    }
    
    configPath := filepath.Join(tempDir, "h3s.yaml")
    err = config.SaveToFile(configPath)
    
    // Assertions
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }
    
    // Verify file exists and contains expected content
    content, err := os.ReadFile(configPath)
    if err != nil {
        t.Fatalf("Failed to read config file: %v", err)
    }
    
    if !strings.Contains(string(content), "test-cluster") {
        t.Errorf("Config file does not contain expected content")
    }
}
```

### Test Data Management

H3S uses the following approaches for test data:

1. **Temporary files and directories**: Use `os.MkdirTemp()` and `os.CreateTemp()` for test files
2. **In-memory data**: Use variables and structs for test data when possible
3. **Test fixtures**: For complex data, use fixtures in the `testdata` directory
4. **Generated test data**: Use helper functions to generate test data programmatically

Example of test data management:

```go
// Create a temporary config file for testing
func createTestConfig(t *testing.T) (string, func()) {
    tempDir, err := os.MkdirTemp("", "h3s-test-*")
    if err != nil {
        t.Fatalf("Failed to create temp directory: %v", err)
    }
    
    configPath := filepath.Join(tempDir, "h3s.yaml")
    config := []byte(`
name: test-cluster
k3s_version: v1.28.2+k3s1
# ... other required fields
`)
    
    if err := os.WriteFile(configPath, config, 0644); err != nil {
        t.Fatalf("Failed to write config file: %v", err)
    }
    
    cleanup := func() {
        os.RemoveAll(tempDir)
    }
    
    return configPath, cleanup
}
```

### Mocking Approach

H3S uses several approaches for mocking:

1. **Interface-based mocking**: Define interfaces and create mock implementations
2. **HTTP server mocking**: Use `httptest.Server` to mock HTTP APIs
3. **Function variable replacement**: Replace function variables with mock implementations
4. **Dependency injection**: Pass mock dependencies to functions and methods

The project includes several mock implementations:

- `mockhetzner`: Mock implementation of the Hetzner Cloud API
- `cmd/testutils`: Utilities for testing CLI commands
- `dependencies.MockDependencies`: Mock implementation of command dependencies

Example of using mocks:

```go
func TestCreateCluster(t *testing.T) {
    // Create a mock server
    mockServer := mockhetzner.NewHetznerMockScenario("/v1", "success")
    defer mockServer.Close()
    
    // Create a mock cluster with the mock server URL
    ctx := &cluster.Cluster{
        Config: &config.Config{
            // ... config fields
            HetznerAPIEndpoint: mockServer.Server.URL,
        },
        Credentials: &credentials.ProjectCredentials{
            // ... credential fields
        },
    }
    
    // Test the function with the mock
    err := Create(ctx)
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }
    
    // Verify the mock was called correctly
    // ... assertions on mock state
}
```

### Coverage Expectations

H3S aims for high test coverage to ensure code quality and reliability:

- **Core business logic**: 80%+ coverage
- **CLI commands**: 70%+ coverage
- **Utility functions**: 60%+ coverage

Coverage is tracked using Codecov and reported in the README.

To check coverage locally:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Pull Request Process

1. Ensure your code passes all tests and linting
2. Update documentation as needed
3. Add or update tests to cover your changes
4. Submit a pull request with a clear description of the changes
5. Address any feedback from reviewers

## Style Guidelines

H3S follows the standard Go style guidelines:

- Use `gofmt` or `go fmt` to format your code
- Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use meaningful variable and function names
- Write clear comments and documentation

## Documentation

- Update the README.md file as needed
- Add godoc comments to exported functions and types
- Update architecture diagrams if you make significant changes
- Add examples for new features
