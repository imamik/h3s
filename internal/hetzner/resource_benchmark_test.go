package hetzner

import (
	"h3s/internal/hetzner/dns"
	"h3s/internal/hetzner/loadbalancers"
	"h3s/internal/hetzner/microos"
	"h3s/internal/hetzner/network"
	"h3s/internal/hetzner/pool"
	"h3s/internal/hetzner/sshkey"
	"testing"
)

// BenchmarkResourceCreation benchmarks the creation of individual resources
func BenchmarkResourceCreation(b *testing.B) {
	// Skip in short mode
	if testing.Short() {
		b.Skip("skipping in short mode")
	}

	// Run sub-benchmarks for each resource type
	b.Run("SSHKey", benchmarkSSHKeyCreation)
	b.Run("Network", benchmarkNetworkCreation)
	b.Run("MicroOS", benchmarkMicroOSCreation)
	b.Run("LoadBalancer", benchmarkLoadBalancerCreation)
	b.Run("Pool", benchmarkPoolCreation)
	b.Run("DNS", benchmarkDNSCreation)
}

// BenchmarkResourceDeletion benchmarks the deletion of individual resources
func BenchmarkResourceDeletion(b *testing.B) {
	// Skip in short mode
	if testing.Short() {
		b.Skip("skipping in short mode")
	}

	// Run sub-benchmarks for each resource type
	b.Run("SSHKey", benchmarkSSHKeyDeletion)
	b.Run("Network", benchmarkNetworkDeletion)
	b.Run("LoadBalancer", benchmarkLoadBalancerDeletion)
	b.Run("Pool", benchmarkPoolDeletion)
	b.Run("DNS", benchmarkDNSDeletion)
}

// Individual resource creation benchmarks

func benchmarkSSHKeyCreation(b *testing.B) {
	// Create a mock cluster
	ctx := createMockCluster(b)

	// Reset the timer before the actual benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create the SSH key
		_, err := sshkey.Create(ctx)
		if err != nil {
			b.Fatalf("failed to create SSH key: %v", err)
		}

		// Reset the mock state between iterations
		b.StopTimer()
		ctx = createMockCluster(b)
		b.StartTimer()
	}
}

func benchmarkNetworkCreation(b *testing.B) {
	// Create a mock cluster
	ctx := createMockCluster(b)

	// Reset the timer before the actual benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create the network
		_, err := network.Create(ctx)
		if err != nil {
			b.Fatalf("failed to create network: %v", err)
		}

		// Reset the mock state between iterations
		b.StopTimer()
		ctx = createMockCluster(b)
		b.StartTimer()
	}
}

func benchmarkMicroOSCreation(b *testing.B) {
	// Create a mock cluster
	ctx := createMockCluster(b)

	// Create SSH key first
	sshKey, err := sshkey.Create(ctx)
	if err != nil {
		b.Fatalf("failed to create SSH key: %v", err)
	}

	// Reset the timer before the actual benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create MicroOS images
		_, err := microos.Create(ctx, sshKey)
		if err != nil {
			b.Fatalf("failed to create MicroOS images: %v", err)
		}

		// Reset the mock state between iterations
		b.StopTimer()
		ctx = createMockCluster(b)
		b.StartTimer()
	}
}

func benchmarkLoadBalancerCreation(b *testing.B) {
	// Create a mock cluster
	ctx := createMockCluster(b)

	// Create network first
	net, err := network.Create(ctx)
	if err != nil {
		b.Fatalf("failed to create network: %v", err)
	}

	// Reset the timer before the actual benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create load balancer
		_, err := loadbalancers.Create(ctx, net)
		if err != nil {
			b.Fatalf("failed to create load balancer: %v", err)
		}

		// Reset the mock state between iterations
		b.StopTimer()
		ctx = createMockCluster(b)
		net, err = network.Create(ctx)
		if err != nil {
			b.Fatalf("failed to create network: %v", err)
		}
		b.StartTimer()
	}
}

func benchmarkPoolCreation(b *testing.B) {
	// Skip this benchmark as it requires more complex setup
	b.Skip("Pool creation benchmark requires more complex setup")

	// This benchmark is skipped, so we don't need the rest of the code
}

func benchmarkDNSCreation(b *testing.B) {
	// Create a mock cluster
	ctx := createMockCluster(b)

	// Create prerequisites
	net, err := network.Create(ctx)
	if err != nil {
		b.Fatalf("failed to create network: %v", err)
	}

	lb, err := loadbalancers.Create(ctx, net)
	if err != nil {
		b.Fatalf("failed to create load balancer: %v", err)
	}

	// Reset the timer before the actual benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create DNS records
		_, err := dns.Create(ctx, lb)
		if err != nil {
			b.Fatalf("failed to create DNS records: %v", err)
		}

		// Reset the mock state between iterations
		b.StopTimer()
		ctx = createMockCluster(b)
		net, _ = network.Create(ctx)
		lb, _ = loadbalancers.Create(ctx, net)
		b.StartTimer()
	}
}

// Individual resource deletion benchmarks

func benchmarkSSHKeyDeletion(b *testing.B) {
	// Create a mock cluster
	ctx := createMockCluster(b)

	// Reset the timer before the actual benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create the SSH key first
		b.StopTimer()
		_, err := sshkey.Create(ctx)
		if err != nil {
			b.Fatalf("failed to create SSH key: %v", err)
		}
		b.StartTimer()

		// Delete the SSH key
		err = sshkey.Delete(ctx)
		if err != nil {
			b.Fatalf("failed to delete SSH key: %v", err)
		}
	}
}

func benchmarkNetworkDeletion(b *testing.B) {
	// Create a mock cluster
	ctx := createMockCluster(b)

	// Reset the timer before the actual benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create the network first
		b.StopTimer()
		_, err := network.Create(ctx)
		if err != nil {
			b.Fatalf("failed to create network: %v", err)
		}
		b.StartTimer()

		// Delete the network
		err = network.Delete(ctx)
		if err != nil {
			b.Fatalf("failed to delete network: %v", err)
		}
	}
}

func benchmarkLoadBalancerDeletion(b *testing.B) {
	// Create a mock cluster
	ctx := createMockCluster(b)

	// Create network first
	net, err := network.Create(ctx)
	if err != nil {
		b.Fatalf("failed to create network: %v", err)
	}

	// Reset the timer before the actual benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create load balancer first
		b.StopTimer()
		_, err := loadbalancers.Create(ctx, net)
		if err != nil {
			b.Fatalf("failed to create load balancer: %v", err)
		}
		b.StartTimer()

		// Delete the load balancer
		err = loadbalancers.Delete(ctx)
		if err != nil {
			b.Fatalf("failed to delete load balancer: %v", err)
		}
	}
}

func benchmarkPoolDeletion(b *testing.B) {
	// Create a mock cluster
	ctx := createMockCluster(b)

	// Skip prerequisites setup as we're just testing the Delete function

	// Reset the timer before the actual benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Skip pool creation as it requires more complex setup
		// This is a placeholder for the actual pool deletion benchmark

		// Delete the pools
		err := pool.Delete(ctx)
		if err != nil {
			b.Fatalf("failed to delete pools: %v", err)
		}
	}
}

func benchmarkDNSDeletion(b *testing.B) {
	// Create a mock cluster
	ctx := createMockCluster(b)

	// Create prerequisites
	net, err := network.Create(ctx)
	if err != nil {
		b.Fatalf("failed to create network: %v", err)
	}

	lb, err := loadbalancers.Create(ctx, net)
	if err != nil {
		b.Fatalf("failed to create load balancer: %v", err)
	}

	// Reset the timer before the actual benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create DNS records first
		b.StopTimer()
		_, err := dns.Create(ctx, lb)
		if err != nil {
			b.Fatalf("failed to create DNS records: %v", err)
		}
		b.StartTimer()

		// Delete the DNS records
		err = dns.Delete(ctx)
		if err != nil {
			b.Fatalf("failed to delete DNS records: %v", err)
		}
	}
}
