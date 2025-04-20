package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"h3s/internal/hetzner/mockhetzner"
)

func TestServerAPI_Success(t *testing.T) {
	mock := mockhetzner.NewHetznerMockScenario("/servers", "success")
	defer mock.Close()
	resp, err := http.Get(mock.Server.URL + "/servers")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
}

func TestServerAPI_Error(t *testing.T) {
	mock := mockhetzner.NewHetznerMockScenario("/servers", "error")
	defer mock.Close()
	resp, err := http.Get(mock.Server.URL + "/servers")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 500, resp.StatusCode)
}

func TestServerAPI_RateLimit(t *testing.T) {
	mock := mockhetzner.NewHetznerMockScenario("/servers", "ratelimit")
	defer mock.Close()
	resp, err := http.Get(mock.Server.URL + "/servers")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 429, resp.StatusCode)
}

func TestServerAPI_Timeout(t *testing.T) {
	mock := mockhetzner.NewHetznerMockScenario("/servers", "timeout")
	defer mock.Close()
	client := &http.Client{Timeout: time.Second}
	resp, err := client.Get(mock.Server.URL + "/servers")
	if err == nil {
		defer resp.Body.Close()
	}
	assert.Error(t, err)
}
