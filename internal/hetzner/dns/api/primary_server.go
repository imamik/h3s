package api

import (
	"context"
	"fmt"
	"net/http"
)

// PrimaryServer represents a primary server in a DNS zone.
type PrimaryServer struct {
	ID      string `json:"id"`
	ZoneID  string `json:"zone_id"`
	Address string `json:"address"`
	Port    uint16 `json:"port"`
}

// CreatePrimaryServerRequest represents a request to create a primary server.
type CreatePrimaryServerRequest struct {
	ZoneID  string `json:"zone_id"`
	Address string `json:"address"`
	Port    uint16 `json:"port"`
}

// PrimaryServersResponse represents a response to get primary servers.
type PrimaryServersResponse struct {
	PrimaryServers []PrimaryServer `json:"primary_servers"`
}

// PrimaryServerResponse represents a response to get a primary server.
type PrimaryServerResponse struct {
	PrimaryServer PrimaryServer `json:"primary_server"`
}

// GetPrimaryServer reads the current state of a primary server.
func (c *Client) GetPrimaryServer(ctx context.Context, id string) (*PrimaryServer, error) {
	resp, err := c.request(ctx, http.MethodGet, "/api/v1/primary_servers/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting primary server %s: %w", id, err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("error closing response body: %w", cerr)
		}
	}()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, nil
	case http.StatusOK:
		var response *PrimaryServerResponse

		err = readAndParseJSONBody(resp, &response)
		if err != nil {
			return nil, fmt.Errorf("error Reading json response of get primary server %s request: %s", id, err)
		}

		return &response.PrimaryServer, nil
	default:
		return nil, fmt.Errorf("http status %d unhandled", resp.StatusCode)
	}
}

// CreatePrimaryServer creates a primary server.
func (c *Client) CreatePrimaryServer(ctx context.Context, server CreatePrimaryServerRequest) (*PrimaryServer, error) {
	resp, err := c.request(ctx, http.MethodPost, "/api/v1/primary_servers", server)
	if err != nil {
		return nil, fmt.Errorf("error creating primary server %s: %w", server.Address, err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("error closing response body: %w", cerr)
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		var response PrimaryServerResponse

		err = readAndParseJSONBody(resp, &response)
		if err != nil {
			return nil, err
		}

		return &response.PrimaryServer, nil
	default:
		return nil, fmt.Errorf("http status %d unhandled", resp.StatusCode)
	}
}

// UpdatePrimaryServer updates a primary server.
func (c *Client) UpdatePrimaryServer(ctx context.Context, server PrimaryServer) (*PrimaryServer, error) {
	resp, err := c.request(ctx, http.MethodPut, "/api/v1/primary_servers/"+server.ID, server)
	if err != nil {
		return nil, fmt.Errorf("error updating primary server %s: %w", server.ID, err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("error closing response body: %w", cerr)
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		var response PrimaryServerResponse

		err = readAndParseJSONBody(resp, &response)
		if err != nil {
			return nil, err
		}

		return &response.PrimaryServer, nil
	default:
		return nil, fmt.Errorf("http status %d unhandled", resp.StatusCode)
	}
}

// DeletePrimaryServer deletes a primary server.
func (c *Client) DeletePrimaryServer(ctx context.Context, id string) error {
	resp, err := c.request(ctx, http.MethodDelete, "/api/v1/primary_servers/"+id, nil)
	if err != nil {
		return fmt.Errorf("error deleting primary server %s: %w", id, err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("error closing response body: %w", cerr)
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	default:
		return fmt.Errorf("http status %d unhandled", resp.StatusCode)
	}
}
