package k8s

import (
	"context"
	"errors"
)

// MockClient simulates Kubernetes resource management and manifest deployment
// Set ShouldFail to true to simulate failure.
type MockClient struct {
	ShouldFail bool
}

func (m *MockClient) ApplyManifest(ctx context.Context, manifest string) error {
	if m.ShouldFail {
		return errors.New("mock k8s apply failed")
	}
	return nil
}

func (m *MockClient) CreateResource(ctx context.Context, kind, name string) error {
	if m.ShouldFail {
		return errors.New("mock k8s create failed")
	}
	return nil
}

func (m *MockClient) DeleteResource(ctx context.Context, kind, name string) error {
	if m.ShouldFail {
		return errors.New("mock k8s delete failed")
	}
	return nil
}
