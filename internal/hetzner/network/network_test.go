package network

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"h3s/internal/hetzner/mockhetzner"
)

func TestNetworkAPI_Success(t *testing.T) {
	mock := mockhetzner.NewHetznerMockScenario("/networks", "success")
	defer mock.Close()
	resp, err := http.Get(mock.Server.URL + "/networks")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
}

func TestNetworkAPI_Error(t *testing.T) {
	mock := mockhetzner.NewHetznerMockScenario("/networks", "error")
	defer mock.Close()
	resp, err := http.Get(mock.Server.URL + "/networks")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 500, resp.StatusCode)
}

func TestNetworkAPI_RateLimit(t *testing.T) {
	mock := mockhetzner.NewHetznerMockScenario("/networks", "ratelimit")
	defer mock.Close()
	resp, err := http.Get(mock.Server.URL + "/networks")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 429, resp.StatusCode)
}

func TestNetworkAPI_Timeout(t *testing.T) {
	mock := mockhetzner.NewHetznerMockScenario("/networks", "timeout")
	defer mock.Close()
	client := &http.Client{Timeout: time.Second}
	resp, err := client.Get(mock.Server.URL + "/networks")
	if err == nil {
		defer resp.Body.Close()
	}
	assert.Error(t, err)
}
