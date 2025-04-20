package pool

import (
	"h3s/internal/cluster"
	"h3s/internal/config"
	"testing"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Minimal stub for cluster.Cluster with safe fields for CreatePools
func newTestCluster() *cluster.Cluster {
	return &cluster.Cluster{
		Config: &config.Config{
			SSHKeyPaths: config.SSHKeyPaths{
				PrivateKeyPath: "/dev/null",
				PublicKeyPath:  "/dev/null",
			},
			NetworkZone: "eu-central",
			K3sVersion:  "v1.0.0",
			Name:        "mock-cluster",
			Domain:      "example.com",
			WorkerPools: []config.NodePool{{
				Instance: "cpx11",
				Location: "nbg1",
				Name:     "worker",
				Nodes:    1,
			}},
			CertManager: config.CertManager{
				Email:      "test@example.com",
				Production: false,
			},
			ControlPlane: config.ControlPlane{
				Pool: config.NodePool{
					Instance: "cpx11",
					Location: "nbg1",
					Name:     "control-plane",
					Nodes:    1,
				},
				AsWorkerPool: false,
			},
		},
		CloudClient: hcloud.NewClient(), // Dummy client
		Context:     nil,
	}
}

func TestPoolCreate_Error(t *testing.T) {
	ctx := newTestCluster()
	_, err := CreatePools(ctx, nil, nil, nil)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestPoolCreate_Success(t *testing.T) {
	t.Skip("cannot run success test without fully mocked dependencies for SSHKey, Network, images")
}
