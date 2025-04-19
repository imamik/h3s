// mockserver.go - Mock HTTP server for API responses
package mockhetzner

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
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

func NewMockServer(handler http.HandlerFunc, config MockServerConfig) *MockServer {
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
		case "timeout":
			time.Sleep(2 * time.Second)
			return
		case "ratelimit":
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":"rate limit exceeded"}`))
			return
		case "error":
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"internal error"}`))
			return
		}
		// Ignore Authorization header for happy path
		switch {
		case r.URL.Path == "/v1/servers":
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				w.Write(mockServers)
			case http.MethodPost:
				w.WriteHeader(201)
				w.Write([]byte(`{"server":{"id":1,"name":"mock-server","status":"running","public_net":{"ipv4":{"ip":"1.2.3.4"},"ipv6":{"ip":"::1"}},"private_net":[{"ip":"10.0.0.1"}]}}`))
			case http.MethodDelete:
				w.WriteHeader(204)
				w.Write([]byte(`{}`))
			}
		case len(r.URL.Path) > len("/v1/servers/") && r.URL.Path[:12] == "/v1/servers/":
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				w.Write([]byte(`{"server":{"id":1,"name":"mock-server","status":"running","public_net":{"ipv4":{"ip":"1.2.3.4"},"ipv6":{"ip":"::1"}},"private_net":[{"ip":"10.0.0.1"}]}}`))
			case http.MethodDelete:
				w.WriteHeader(204)
				w.Write([]byte(`{}`))
			}
		case r.URL.Path == "/v1/ssh_keys":
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				w.Write(mockSSHKeys)
			case http.MethodPost:
				w.WriteHeader(201)
				w.Write([]byte(`{"ssh_key":{"id":1,"name":"mock-key","public_key":"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMockKey== mock@mock","fingerprint":"mockfp"}}`))
			case http.MethodDelete:
				w.WriteHeader(204)
				w.Write([]byte(`{}`))
			}
		case len(r.URL.Path) > len("/v1/ssh_keys/") && r.URL.Path[:13] == "/v1/ssh_keys/":
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				w.Write([]byte(`{"ssh_key":{"id":1,"name":"mock-key","public_key":"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMockKey== mock@mock","fingerprint":"mockfp"}}`))
			case http.MethodDelete:
				w.WriteHeader(204)
				w.Write([]byte(`{}`))
			}
		case r.URL.Path == "/v1/networks":
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				w.Write(mockNetworks)
			case http.MethodPost:
				w.WriteHeader(201)
				w.Write([]byte(`{"network":{"id":1,"name":"mock-network","ip_range":"10.0.0.0/16"}}`))
			case http.MethodDelete:
				w.WriteHeader(204)
				w.Write([]byte(`{}`))
			}
		case len(r.URL.Path) > len("/v1/networks/") && r.URL.Path[:14] == "/v1/networks/":
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				w.Write([]byte(`{"network":{"id":1,"name":"mock-network","ip_range":"10.0.0.0/16"}}`))
			case http.MethodDelete:
				w.WriteHeader(204)
				w.Write([]byte(`{}`))
			}
		case r.URL.Path == "/v1/load_balancers":
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				w.Write(mockLoadBalancers)
			case http.MethodPost:
				w.WriteHeader(201)
				w.Write([]byte(`{"load_balancer":{"id":1,"name":"mock-lb","public_net":{"ipv4":{"ip":"2.2.2.2"},"ipv6":{"ip":"::2"}}}}`))
			case http.MethodDelete:
				w.WriteHeader(204)
				w.Write([]byte(`{}`))
			}
		case len(r.URL.Path) > len("/v1/load_balancers/") && r.URL.Path[:19] == "/v1/load_balancers/":
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				w.Write([]byte(`{"load_balancer":{"id":1,"name":"mock-lb","public_net":{"ipv4":{"ip":"2.2.2.2"},"ipv6":{"ip":"::2"}}}}`))
			case http.MethodDelete:
				w.WriteHeader(204)
				w.Write([]byte(`{}`))
			}
		default:
			w.WriteHeader(200)
			w.Write([]byte(`{}`))
		}
	}))
	ms.Server = ts
	return ms
}

func (m *MockServer) Close() {
	m.Server.Close()
}

// Helper to easily set up a scenario
func NewHetznerMockScenario(endpoint string, mode string) *MockServer {
	cfg := MockServerConfig{ResponseCode: 200, ResponseBody: `{"result":"ok"}`}
	switch mode {
	case "success":
		cfg.ResponseCode = 200
		cfg.ResponseBody = `{"result":"ok"}`
	case "error":
		cfg.ErrorMode = "error"
	case "ratelimit":
		cfg.ErrorMode = "ratelimit"
	case "timeout":
		cfg.ErrorMode = "timeout"
		cfg.Delay = time.Second * 2
	}
	return NewMockServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), cfg)
}
