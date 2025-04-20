// mockserver.go - Mock HTTP server for API responses
package mockhetzner

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

// Error modes for the mock server
const (
	timeoutMode   = "timeout"
	ratelimitMode = "ratelimit"
	errorMode     = "error"
)

type MockServer struct {
	Server   *httptest.Server
	Requests []*http.Request
	mu       sync.Mutex
	Config   MockServerConfig
}

type MockServerConfig struct {
	ResponseCode int
	ResponseBody string
	Delay        time.Duration
	RateLimit    int
	ErrorMode    string // "", "error", "timeout", "ratelimit"
}

//nolint:gocyclo // Complexity acceptable for mock server setup
func NewMockServer(_ http.HandlerFunc, config MockServerConfig) *MockServer {
	ms := &MockServer{Config: config}
	// In-memory state for servers and ssh keys
	var (
		mockServers       = []byte(`{"servers":[{"id":1,"name":"mock-server","status":"running","public_net":{"ipv4":{"ip":"1.2.3.4"},"ipv6":{"ip":"::1"}},"private_net":[{"ip":"10.0.0.1"}]}]}`)
		mockSSHKeys       = []byte(`{"ssh_keys":[{"id":1,"name":"mock-key","public_key":"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMockKey== mock@mock","fingerprint":"mockfp"}]}`)
		mockNetworks      = []byte(`{"networks":[{"id":1,"name":"mock-network","ip_range":"10.0.0.0/16"}]}`)
		mockLoadBalancers = []byte(`{"load_balancers":[{"id":1,"name":"mock-lb","public_net":{"ipv4":{"ip":"2.2.2.2"},"ipv6":{"ip":"::2"}}}]}`)
	)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ms.mu.Lock()
		ms.Requests = append(ms.Requests, r)
		ms.mu.Unlock()
		// Log request method, path, and headers for debugging
		println("MOCK REQUEST:", r.Method, r.URL.Path)
		for name, values := range r.Header {
			for _, v := range values {
				println("  ", name, ":", v)
			}
		}
		if ms.Config.Delay > 0 {
			time.Sleep(ms.Config.Delay)
		}
		switch ms.Config.ErrorMode {
		case timeoutMode:
			time.Sleep(2 * time.Second)
			return
		case ratelimitMode:
			w.WriteHeader(http.StatusTooManyRequests)
			safeWrite(w, []byte(`{"error":"rate limit exceeded"}`))
			return
		case errorMode:
			w.WriteHeader(http.StatusInternalServerError)
			safeWrite(w, []byte(`{"error":"internal error"}`))
			return
		}
		// Ignore Authorization header for happy path
		switch {
		case r.URL.Path == "/v1/servers":
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				safeWrite(w, mockServers)
			case http.MethodPost:
				w.WriteHeader(201)
				safeWrite(w, []byte(`{"server":{"id":1,"name":"mock-server","status":"running","public_net":{"ipv4":{"ip":"1.2.3.4"},"ipv6":{"ip":"::1"}},"private_net":[{"ip":"10.0.0.1"}]}}`))
			case http.MethodDelete:
				w.WriteHeader(204)
				safeWrite(w, []byte(`{}`))
			}
		case len(r.URL.Path) > len("/v1/servers/") && r.URL.Path[:12] == "/v1/servers/":
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				safeWrite(w, []byte(`{"server":{"id":1,"name":"mock-server","status":"running","public_net":{"ipv4":{"ip":"1.2.3.4"},"ipv6":{"ip":"::1"}},"private_net":[{"ip":"10.0.0.1"}]}}`))
			case http.MethodDelete:
				w.WriteHeader(204)
				safeWrite(w, []byte(`{}`))
			}
		case r.URL.Path == "/v1/ssh_keys":
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				safeWrite(w, mockSSHKeys)
			case http.MethodPost:
				w.WriteHeader(201)
				safeWrite(w, []byte(`{"ssh_key":{"id":1,"name":"mock-key","public_key":"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMockKey== mock@mock","fingerprint":"mockfp"}}`))
			case http.MethodDelete:
				w.WriteHeader(204)
				safeWrite(w, []byte(`{}`))
			}
		case len(r.URL.Path) > len("/v1/ssh_keys/") && r.URL.Path[:13] == "/v1/ssh_keys/":
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				safeWrite(w, []byte(`{"ssh_key":{"id":1,"name":"mock-key","public_key":"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMockKey== mock@mock","fingerprint":"mockfp"}}`))
			case http.MethodDelete:
				w.WriteHeader(204)
				safeWrite(w, []byte(`{}`))
			}
		case r.URL.Path == "/v1/networks":
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				safeWrite(w, mockNetworks)
			case http.MethodPost:
				w.WriteHeader(201)
				safeWrite(w, []byte(`{"network":{"id":1,"name":"mock-network","ip_range":"10.0.0.0/16"}}`))
			case http.MethodDelete:
				w.WriteHeader(204)
				safeWrite(w, []byte(`{}`))
			}
		case len(r.URL.Path) > len("/v1/networks/") && r.URL.Path[:14] == "/v1/networks/":
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				safeWrite(w, []byte(`{"network":{"id":1,"name":"mock-network","ip_range":"10.0.0.0/16"}}`))
			case http.MethodDelete:
				w.WriteHeader(204)
				safeWrite(w, []byte(`{}`))
			}
		case r.URL.Path == "/v1/load_balancers":
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				safeWrite(w, mockLoadBalancers)
			case http.MethodPost:
				w.WriteHeader(201)
				safeWrite(w, []byte(`{"load_balancer":{"id":1,"name":"mock-lb","public_net":{"ipv4":{"ip":"2.2.2.2"},"ipv6":{"ip":"::2"}}}}`))
			case http.MethodDelete:
				w.WriteHeader(204)
				safeWrite(w, []byte(`{}`))
			}
		case len(r.URL.Path) > len("/v1/load_balancers/") && r.URL.Path[:19] == "/v1/load_balancers/":
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				safeWrite(w, []byte(`{"load_balancer":{"id":1,"name":"mock-lb","public_net":{"ipv4":{"ip":"2.2.2.2"},"ipv6":{"ip":"::2"}}}}`))
			case http.MethodDelete:
				w.WriteHeader(204)
				safeWrite(w, []byte(`{}`))
			}
		default:
			w.WriteHeader(200)
			safeWrite(w, []byte(`{}`))
		}
	}))
	ms.Server = ts
	return ms
}

func (m *MockServer) Close() {
	m.Server.Close()
}

// safeWrite is a helper function to safely write to an http.ResponseWriter
func safeWrite(w http.ResponseWriter, data []byte) {
	if _, err := w.Write(data); err != nil {
		println("Error writing response:", err)
	}
}

// NewHetznerMockScenario creates a new mock server with predefined scenarios.
//
//nolint:gocyclo // Complex setup logic is acceptable for mock server initialization
func NewHetznerMockScenario(
	_ string,
	mode string,
) *MockServer {
	cfg := MockServerConfig{ResponseCode: 200, ResponseBody: `{"result":"ok"}`}
	switch mode {
	case "success":
		cfg.ResponseCode = 200
		cfg.ResponseBody = `{"result":"ok"}`
	case errorMode:
		cfg.ErrorMode = errorMode
	case ratelimitMode:
		cfg.ErrorMode = ratelimitMode
	case timeoutMode:
		cfg.ErrorMode = timeoutMode
		cfg.Delay = time.Second * 2
	}
	return NewMockServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}), cfg)
}
