package k8s

import (
	"context"
	"testing"
)

func TestMockClient_ApplyManifest_Success(t *testing.T) {
	m := &MockClient{ShouldFail: false}
	if err := m.ApplyManifest(context.Background(), "apiVersion: v1\nkind: Pod"); err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

func TestMockClient_ApplyManifest_Fail(t *testing.T) {
	m := &MockClient{ShouldFail: true}
	if err := m.ApplyManifest(context.Background(), "apiVersion: v1\nkind: Pod"); err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestMockClient_CreateResource_Success(t *testing.T) {
	m := &MockClient{ShouldFail: false}
	if err := m.CreateResource(context.Background(), "Pod", "foo"); err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

func TestMockClient_CreateResource_Fail(t *testing.T) {
	m := &MockClient{ShouldFail: true}
	if err := m.CreateResource(context.Background(), "Pod", "foo"); err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestMockClient_DeleteResource_Success(t *testing.T) {
	m := &MockClient{ShouldFail: false}
	if err := m.DeleteResource(context.Background(), "Pod", "foo"); err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

func TestMockClient_DeleteResource_Fail(t *testing.T) {
	m := &MockClient{ShouldFail: true}
	if err := m.DeleteResource(context.Background(), "Pod", "foo"); err == nil {
		t.Errorf("expected error, got nil")
	}
}
