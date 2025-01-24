// Package api contains the Hetzner DNS API client
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Client for the Hetzner DNS API.
type Client struct {
	httpClient  *http.Client
	endPoint    *url.URL
	apiToken    string
	userAgent   string
	requestLock sync.Mutex
}

// New creates a new API Client using a given api token.
func New(apiEndpoint, apiToken string, roundTripper http.RoundTripper) (*Client, error) {
	endPoint, err := url.Parse(apiEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error parsing API endpoint URL: %w", err)
	}

	httpClient := &http.Client{
		Transport: roundTripper,
	}

	client := &Client{
		apiToken:   apiToken,
		endPoint:   endPoint,
		httpClient: httpClient,
	}

	return client, nil
}

// SetUserAgent sets the user agent for the client.
func (c *Client) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
}

// request sends a request to the Hetzner DNS API and returns the response.
func (c *Client) request(ctx context.Context, method, path string, bodyJSON any) (*http.Response, error) {
	uri := c.endPoint.String() + path
	tflog.Debug(ctx, fmt.Sprintf("HTTP request to API %s %s", method, uri))

	reqBody, err := c.prepareRequestBody(bodyJSON)
	if err != nil {
		return nil, err
	}

	req, err := c.buildRequest(ctx, method, uri, reqBody)
	if err != nil {
		return nil, err
	}

	c.lockRequest(method)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	return c.handleResponse(resp)
}

func (c *Client) prepareRequestBody(bodyJSON any) ([]byte, error) {
	if bodyJSON == nil {
		return nil, nil
	}

	reqBody, err := json.Marshal(bodyJSON)
	if err != nil {
		return nil, fmt.Errorf("error serializing JSON body %s", err)
	}

	return reqBody, nil
}

func (c *Client) buildRequest(ctx context.Context, method, uri string, reqBody []byte) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, uri, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}

	req.Header.Set("Auth-API-Token", c.apiToken)
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	return req, nil
}

func (c *Client) lockRequest(method string) {
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete {
		c.requestLock.Lock()
		defer c.requestLock.Unlock()
	}
}

func (c *Client) handleResponse(resp *http.Response) (*http.Response, error) {
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		unauthorizedError, err := parseUnauthorizedError(resp)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("API returned HTTP 401 Unauthorized error with message: '%s'. "+
			"Double check your API key is still valid", unauthorizedError.Message)
	case http.StatusUnprocessableEntity:
		unprocessableEntityError, err := parseUnprocessableEntityError(resp)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("API returned HTTP 422 Unprocessable Entity error with message: '%s'", unprocessableEntityError.Error.Message)
	default:
		return resp, nil
	}
}
