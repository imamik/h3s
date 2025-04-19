package hetzner

import (
	"context"
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/config/credentials"
	"h3s/internal/hetzner/dns/api"
	"h3s/internal/hetzner/mockhetzner"
	"net/http"
	"testing"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// createMockCluster creates a mock cluster for benchmarking
func createMockCluster(b *testing.B) *cluster.Cluster {
	// Create a mock server
	mock := mockhetzner.NewHetznerMockScenario("/v1", "success")

	// Create a minimal config
	conf := &config.Config{
		Name: "benchmark-cluster",
		ControlPlane: config.ControlPlane{
			Pool: config.NodePool{
				Instance: "cx31",
				Location: "nbg1",
				Name:     "cp01",
				Nodes:    1,
			},
		},
		WorkerPools: []config.NodePool{
			{
				Instance: "cx31",
				Location: "nbg1",
				Name:     "wp1",
				Nodes:    1,
			},
		},
		NetworkZone: "nbg1",
		Domain:      "benchmark.example.com",
	}

	// Create mock credentials
	creds := &credentials.ProjectCredentials{
		HCloudToken:     "p1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcde",
		HetznerDNSToken: "1234567890abcdef1234567890abcdef",
		K3sToken:        "k3s1234567890abcdef1234567890abcdef",
	}

	// Create a client using the mock server
	client := hcloud.NewClient(
		hcloud.WithToken(creds.HCloudToken),
		hcloud.WithEndpoint(mock.Server.URL),
	)

	// Create a DNS client
	dnsClient, err := api.New(mock.Server.URL, creds.HetznerDNSToken, http.DefaultTransport)
	if err != nil {
		b.Fatalf("failed to create DNS client: %v", err)
	}

	// Return a new Cluster with all components initialized
	return &cluster.Cluster{
		Config:      conf,
		Credentials: creds,
		CloudClient: client,
		DNSClient:   dnsClient,
		Context:     context.Background(),
	}
}

// BenchmarkCreateDestroy benchmarks the create and destroy operations
func BenchmarkCreateDestroy(b *testing.B) {
	// Skip in short mode
	if testing.Short() {
		b.Skip("skipping in short mode")
	}

	// Create a mock cluster
	ctx := createMockCluster(b)

	// Reset the timer before the actual benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create the cluster
		err := Create(ctx)
		if err != nil {
			b.Fatalf("failed to create cluster: %v", err)
		}

		// Destroy the cluster
		err = Destroy(ctx)
		if err != nil {
			b.Fatalf("failed to destroy cluster: %v", err)
		}
	}
}

// BenchmarkCreateOnly benchmarks just the create operation
func BenchmarkCreateOnly(b *testing.B) {
	// Skip in short mode
	if testing.Short() {
		b.Skip("skipping in short mode")
	}

	// Create a mock cluster
	ctx := createMockCluster(b)

	// Reset the timer before the actual benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create the cluster
		err := Create(ctx)
		if err != nil {
			b.Fatalf("failed to create cluster: %v", err)
		}

		// Reset the mock state between iterations
		b.StopTimer()
		ctx = createMockCluster(b)
		b.StartTimer()
	}
}

// BenchmarkDestroyOnly benchmarks just the destroy operation
func BenchmarkDestroyOnly(b *testing.B) {
	// Skip in short mode
	if testing.Short() {
		b.Skip("skipping in short mode")
	}

	// Create a mock cluster
	ctx := createMockCluster(b)

	// Create the cluster once before benchmarking destroy
	err := Create(ctx)
	if err != nil {
		b.Fatalf("failed to create cluster: %v", err)
	}

	// Reset the timer before the actual benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Destroy the cluster
		err = Destroy(ctx)
		if err != nil {
			b.Fatalf("failed to destroy cluster: %v", err)
		}

		// Reset the mock state between iterations
		b.StopTimer()
		ctx = createMockCluster(b)
		err = Create(ctx)
		if err != nil {
			b.Fatalf("failed to create cluster: %v", err)
		}
		b.StartTimer()
	}
}

// BenchmarkParallelCreateDestroy benchmarks parallel create and destroy operations
func BenchmarkParallelCreateDestroy(b *testing.B) {
	// Skip in short mode
	if testing.Short() {
		b.Skip("skipping in short mode")
	}

	// Reset the timer before the actual benchmark
	b.ResetTimer()

	// Run the benchmark in parallel
	b.RunParallel(func(pb *testing.PB) {
		// Create a mock cluster for each goroutine
		ctx := createMockCluster(b)

		for pb.Next() {
			// Create the cluster
			err := Create(ctx)
			if err != nil {
				b.Fatalf("failed to create cluster: %v", err)
			}

			// Destroy the cluster
			err = Destroy(ctx)
			if err != nil {
				b.Fatalf("failed to destroy cluster: %v", err)
			}

			// Reset the mock state between iterations
			ctx = createMockCluster(b)
		}
	})
}
