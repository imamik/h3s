package config

import (
	"os"
	"path/filepath"
	"testing"

	"h3s/internal/utils/file"
	"h3s/internal/validation"
)

// BenchmarkConfigLoad benchmarks the loading and parsing of configuration files
func BenchmarkConfigLoad(b *testing.B) {
	// Create a temporary config file for benchmarking
	content := []byte(`
ssh_key_paths:
  private_key_path: /tmp/id_rsa
  public_key_path: /tmp/id_rsa.pub
network_zone: nbg1
k3s_version: v1.28.0
name: benchcluster
domain: example.com
cert_manager:
  email: test@example.com
  production: false
worker_pools:
  - instance: cx31
    location: nbg1
    name: workerpool1
    nodes: 1
control_plane:
  pool:
    instance: cx31
    location: nbg1
    name: cp01
    nodes: 1
  as_worker_pool: false
`)
	tempfile, err := os.CreateTemp("", "bench-config-*.yaml")
	if err != nil {
		b.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempfile.Name())
	if _, err := tempfile.Write(content); err != nil {
		b.Fatalf("failed to write to temp file: %v", err)
	}
	tempfile.Close()
	os.Setenv("H3S_CONFIG", tempfile.Name())
	defer os.Unsetenv("H3S_CONFIG")

	// Reset the timer before the actual benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		_, err := Load()
		if err != nil {
			b.Fatalf("failed to load config: %v", err)
		}
	}
}

// BenchmarkConfigLoadLarge benchmarks loading and parsing a large configuration file
func BenchmarkConfigLoadLarge(b *testing.B) {
	// Create a temporary config file with many worker pools for benchmarking
	baseContent := `
ssh_key_paths:
  private_key_path: /tmp/id_rsa
  public_key_path: /tmp/id_rsa.pub
network_zone: nbg1
k3s_version: v1.28.0
name: benchcluster
domain: example.com
cert_manager:
  email: test@example.com
  production: false
control_plane:
  pool:
    instance: cx31
    location: nbg1
    name: cp01
    nodes: 1
  as_worker_pool: false
worker_pools:
`
	// Add 50 worker pools to make the config large
	workerPoolContent := ""
	for i := 0; i < 50; i++ {
		workerPoolContent += `  - instance: cx31
    location: nbg1
    name: workerpool` + string(rune('a'+i%26)) + `
    nodes: 1
`
	}
	content := []byte(baseContent + workerPoolContent)

	tempfile, err := os.CreateTemp("", "bench-config-large-*.yaml")
	if err != nil {
		b.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempfile.Name())
	if _, err := tempfile.Write(content); err != nil {
		b.Fatalf("failed to write to temp file: %v", err)
	}
	tempfile.Close()
	os.Setenv("H3S_CONFIG", tempfile.Name())
	defer os.Unsetenv("H3S_CONFIG")

	// Reset the timer before the actual benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		_, err := Load()
		if err != nil {
			b.Fatalf("failed to load config: %v", err)
		}
	}
}

// BenchmarkGetArchitectures benchmarks the GetArchitectures function
func BenchmarkGetArchitectures(b *testing.B) {
	// Create a config with mixed architectures
	config := &Config{
		ControlPlane: ControlPlane{
			Pool: NodePool{
				Instance: "cx31", // x86
				Location: "nbg1",
				Name:     "cp01",
				Nodes:    1,
			},
		},
		WorkerPools: []NodePool{
			{
				Instance: "cx31", // x86
				Location: "nbg1",
				Name:     "wp1",
				Nodes:    1,
			},
			{
				Instance: "cax11", // ARM
				Location: "nbg1",
				Name:     "wp2",
				Nodes:    1,
			},
		},
	}

	// Reset the timer before the actual benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		_ = GetArchitectures(config)
	}
}

// BenchmarkConfigValidation benchmarks the validation of configuration
func BenchmarkConfigValidation(b *testing.B) {
	// Create a valid config for benchmarking
	config := &Config{
		SSHKeyPaths: SSHKeyPaths{
			PrivateKeyPath: "/tmp/id_rsa",
			PublicKeyPath:  "/tmp/id_rsa.pub",
		},
		NetworkZone: "nbg1",
		K3sVersion:  "v1.28.0",
		Name:        "benchcluster",
		Domain:      "example.com",
		CertManager: CertManager{
			Email:      "test@example.com",
			Production: false,
		},
		ControlPlane: ControlPlane{
			Pool: NodePool{
				Instance: "cx31",
				Location: "nbg1",
				Name:     "cp01",
				Nodes:    1,
			},
			AsWorkerPool: false,
		},
		WorkerPools: []NodePool{
			{
				Instance: "cx31",
				Location: "nbg1",
				Name:     "workerpool1",
				Nodes:    1,
			},
		},
	}

	// Reset the timer before the actual benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		err := validation.ValidateStruct(config)
		if err != nil {
			b.Fatalf("failed to validate config: %v", err)
		}
	}
}

// BenchmarkConfigSave benchmarks saving a configuration to a file
func BenchmarkConfigSave(b *testing.B) {
	// Create a config for benchmarking
	config := &Config{
		SSHKeyPaths: SSHKeyPaths{
			PrivateKeyPath: "/tmp/id_rsa",
			PublicKeyPath:  "/tmp/id_rsa.pub",
		},
		NetworkZone: "nbg1",
		K3sVersion:  "v1.28.0",
		Name:        "benchcluster",
		Domain:      "example.com",
		CertManager: CertManager{
			Email:      "test@example.com",
			Production: false,
		},
		ControlPlane: ControlPlane{
			Pool: NodePool{
				Instance: "cx31",
				Location: "nbg1",
				Name:     "cp01",
				Nodes:    1,
			},
			AsWorkerPool: false,
		},
		WorkerPools: []NodePool{
			{
				Instance: "cx31",
				Location: "nbg1",
				Name:     "wp1",
				Nodes:    1,
			},
		},
	}

	// Create a temporary directory for benchmark files
	tempDir, err := os.MkdirTemp("", "bench-config-save-*")
	if err != nil {
		b.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Reset the timer before the actual benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create a unique filename for each iteration
		filename := filepath.Join(tempDir, "config-"+string(rune('a'+i%26))+".yaml")

		// Save the config to a file
		_, err := file.New(filename).SetYaml(config).Save()
		if err != nil {
			b.Fatalf("failed to save config: %v", err)
		}
	}
}
