# Technical Design Document: Comprehensive and Detailed Testing Strategy for h3s

## Table of Contents
1. Introduction
2. Codebase Mapping and Module Analysis
3. Test Architecture and Principles
4. Detailed Test Plan by Module
5. Test Data, Mocking, and Isolation
6. Automation, CI/CD, and Developer Experience
7. Performance, Security, and Fuzz Testing
8. Documentation and Reporting
9. Risks, Technical Debt, and Continuous Improvement
10. Appendix: Example Test Skeletons

---

## 1. Introduction

This document presents a comprehensive, actionable, and highly detailed technical design for testing the `h3s` CLI tool. The goal is to achieve and sustain 80-90% code coverage, ensure robust quality, and provide a foundation for maintainable and scalable development.

---

## 2. Codebase Mapping and Module Analysis

### Top-Level Structure
- `/cmd`: CLI entrypoints and subcommands (create, destroy, get, install, kubectl, ssh, version)
- `/internal`: Core business logic, organized by domain
  - `cluster`: Cluster lifecycle management
  - `config`: Parsing, validation, and management of configuration files
  - `hetzner`: Hetzner Cloud API interactions (servers, networks, firewalls, etc.)
  - `k3s`: K3s-specific logic (installation, upgrades, etc.)
  - `k8s`: Kubernetes resource management (manifests, deployments)
  - `utils`: Utility functions (file I/O, string manipulation, logging)
  - `validation`: Data and config validation
  - `version`: Version/build info
- `/test/integration`: End-to-end and workflow tests
- `main.go`: CLI entrypoint
- `main_test.go`: Top-level CLI tests
- `.github/`, `Makefile`, `scripts/`: CI/CD, automation, and scripts

### Module Interactions
- CLI commands in `/cmd` invoke business logic in `/internal` modules.
- `/internal/hetzner` and `/internal/k3s` interact with external APIs.
- `/internal/config` and `/internal/validation` ensure input correctness.
- `/test/integration` simulates real-world workflows, using dummy configs and secrets.

### Observed Issues
- Some business logic is tightly coupled to CLI, limiting testability.
- Error handling and logging are not fully standardized.
- Potential for logic duplication (e.g., resource cleanup, API retries).
- Utilities and edge cases are under-tested.
- Integration tests may not fully isolate side effects.

---

## 3. Test Architecture and Principles

### Test Types
- **Unit Tests**: Isolate and test individual functions/methods; mock all dependencies.
- **Integration Tests**: Validate interactions between modules and simulate real workflows.
- **System/E2E Tests**: Run full CLI flows in a sandboxed environment.
- **Negative/Edge Tests**: Intentionally break things (invalid inputs, API failures, etc.).
- **Performance Tests**: Benchmark critical operations.
- **Security/Fuzz Tests**: Probe for vulnerabilities and parser robustness.

### Principles
- Tests must be deterministic, isolated, and repeatable.
- Each test should have a clear purpose and document its intent.
- All external dependencies must be mocked or sandboxed.
- Test data must be managed and cleaned up after each run.

---

## 4. Detailed Test Plan by Module

### `/cmd` (CLI Layer)
- **Test all subcommands**: Argument parsing, flag handling, output (success, error, help), version/help flags.
- **Table-driven tests** for each command, covering all flag/argument combinations.
- **Invalid input tests**: Unknown subcommands, missing required flags, conflicting flags.
- **Output validation**: Check for expected stdout/stderr output and exit codes.
- **Mock business logic**: Use dependency injection or interfaces to replace real logic with mocks.

### `/internal/cluster`
- **Unit tests** for cluster creation, deletion, and update logic:
  - Valid/invalid configs
  - Idempotency (repeated create/delete)
  - Error scenarios (API/network failures, permission errors)
- **Mock Hetzner and K3s APIs**: Simulate all API responses, including errors and rate limits.
- **Test resource cleanup**: Ensure all created resources are deleted on failure.
- **Edge cases**: Large clusters, duplicate names, resource exhaustion.

### `/internal/config` & `/internal/validation`
- **Unit tests** for config parsing and validation:
  - Valid, missing, and malformed YAML
  - Boundary values for all fields
  - Defaulting logic
  - Invalid field types/values
- **Fuzz tests**: Randomized YAML input to catch parser vulnerabilities.

### `/internal/hetzner`
- **Unit tests** for all API wrappers:
  - Success, failure, and retry logic
  - Rate limiting and timeout handling
  - Invalid API credentials
- **Mock all HTTP responses**: Use a test server or mocking library.
- **Coverage of all resource types**: Servers, networks, firewalls, load balancers, etc.

### `/internal/k3s` & `/internal/k8s`
- **Unit tests** for manifest generation, validation, and deployment logic:
  - Valid/invalid manifests
  - Unsupported options
  - Resource limits and quotas
- **Integration tests**: Simulate end-to-end deployment flows.

### `/internal/utils`
- **Unit tests** for all utility functions:
  - File I/O (read/write, permission errors)
  - String manipulation (trimming, splitting, edge cases)
  - Logging (output format, error vs info)
- **Negative tests**: Unreadable files, empty/large inputs.

### `/internal/validation`
- **Unit tests** for all validators:
  - Valid, invalid, and boundary data
  - Error message accuracy

### `/internal/version`
- **Unit tests** for version/build info injection:
  - Default values
  - Values injected via ldflags

### `/test/integration`
- **Full lifecycle tests**: Create, validate, and destroy clusters.
- **Config/secret scenarios**: Valid, missing, malformed, and permission-restricted files.
- **API/network failure simulation**: Use dummy tokens and simulate API downtime.
- **Cleanup checks**: Ensure no residue after test runs.
- **Concurrent operations**: If supported, test parallel cluster operations.

---

## 5. Test Data, Mocking, and Isolation

- Use dedicated test directories for configs, secrets, and temp files.
- Use dummy tokens and credentials for integration tests.
- Clean up all files and resources after each test (defer os.Remove, etc.).
- Use interfaces for all external dependencies (API, file system, network).
- Use Go’s `testing` package and `testify` for assertions and mocks.
- Avoid global state; reset between tests.

---

## 6. Automation, CI/CD, and Developer Experience

- **Makefile targets**: `make test`, `make lint`, `make coverage`, `make integration`
- **Pre-commit hooks**: Format, lint, and run unit tests on commit.
- **CI/CD**: Run all tests and collect coverage on PR and push (GitHub Actions or similar).
- **Coverage enforcement**: Fail builds on regression below threshold.
- **Coverage badge**: Display in README.
- **Automated dependency updates**: Use Dependabot or similar.

---

## 7. Performance, Security, and Fuzz Testing

- **Performance benchmarks**: Use Go’s `testing.B` for critical paths (config parsing, API calls).
- **Load tests**: Simulate large cluster definitions and concurrent operations.
- **Security tests**:
  - Secrets/config leakage (logs, errors)
  - Permission and access control (read-only files, missing keys)
  - Fuzz config and secret parsers for injection vulnerabilities
- **Static analysis**: Run `go vet`, `golangci-lint`, and security linters in CI.

---

## 8. Documentation and Reporting

- **GoDoc**: All exported functions/types must have documentation.
- **Test case documentation**: Each test describes its intent and expected outcome.
- **README**: Document test strategies, how to run tests, and coverage expectations.
- **CONTRIBUTING.md**: Onboarding for new contributors, including testing standards.
- **Coverage reports**: Store and review after each CI run.

---

## 9. Risks, Technical Debt, and Continuous Improvement

- **Global state**: Refactor to minimize; use dependency injection.
- **Error handling**: Standardize with error wrapping/context.
- **Duplicated logic**: Refactor common code (e.g., resource cleanup, API retries).
- **Large functions**: Break up into smaller, testable units.
- **Direct file/network access**: Always abstract behind interfaces.
- **Regular reviews**: Track coverage, technical debt, and address via GitHub issues.

---

## 10. Appendix: Example Test Skeletons

### Table-Driven Unit Test for Config Validation
```go
type configTestCase struct {
    name     string
    input    []byte
    wantErr  bool
}

func TestParseConfig(t *testing.T) {
    cases := []configTestCase{
        {"valid config", []byte(validYAML), false},
        {"missing field", []byte(missingFieldYAML), true},
        {"malformed yaml", []byte(malformedYAML), true},
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            _, err := ParseConfig(tc.input)
            if (err != nil) != tc.wantErr {
                t.Errorf("got error %v, wantErr %v", err, tc.wantErr)
            }
        })
    }
}
```

### Mocked Hetzner API Test
```go
func TestCreateServer(t *testing.T) {
    mockAPI := &MockHetznerAPI{ /* configure responses */ }
    cluster := NewCluster(mockAPI)
    err := cluster.CreateServer("test-server")
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    // Assert API was called as expected
}
```

### CLI Integration Test
```go
func TestCreateClusterCommand(t *testing.T) {
    cmd := exec.Command("go", "run", "../..", "create", "cluster", "--config", "testdata/valid.yaml")
    output, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("command failed: %v, output: %s", err, output)
    }
    if !strings.Contains(string(output), "Cluster created successfully") {
        t.Errorf("unexpected output: %s", output)
    }
}
```

---

# End of Document
