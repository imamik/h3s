package k3s

import "testing"

func TestMockInstaller_Install_Success(t *testing.T) {
	m := &MockInstaller{ShouldFail: false}
	if err := m.Install(); err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

func TestMockInstaller_Install_Fail(t *testing.T) {
	m := &MockInstaller{ShouldFail: true}
	if err := m.Install(); err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestMockInstaller_Configure_Success(t *testing.T) {
	m := &MockInstaller{ShouldFail: false}
	if err := m.Configure(); err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

func TestMockInstaller_Configure_Fail(t *testing.T) {
	m := &MockInstaller{ShouldFail: true}
	if err := m.Configure(); err == nil {
		t.Errorf("expected error, got nil")
	}
}
