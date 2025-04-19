package k3s

import "errors"

// MockInstaller simulates K3s installation/configuration
// Set ShouldFail to true to simulate failure.
type MockInstaller struct {
	ShouldFail bool
}

func (m *MockInstaller) Install() error {
	if m.ShouldFail {
		return errors.New("mock k3s install failed")
	}
	return nil
}

func (m *MockInstaller) Configure() error {
	if m.ShouldFail {
		return errors.New("mock k3s config failed")
	}
	return nil
}
