package api

import (
	"context"
	"fmt"
	"net/http"
)

// Record represents a record in a specific Zone.
type Record struct {
	TTL    *int64 `json:"ttl,omitempty"`
	ZoneID string `json:"zone_id"`
	ID     string `json:"id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Value  string `json:"value"`
}

// HasTTL returns true if a Record has a TTL set and false if TTL is undefined.
func (r *Record) HasTTL() bool {
	return r.TTL != nil
}

// CreateRecordOpts covers all parameters used to create a new DNS record.
type CreateRecordOpts struct {
	TTL    *int64 `json:"ttl,omitempty"`
	ZoneID string `json:"zone_id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Value  string `json:"value"`
}

// CreateRecordRequest represents all data required to create a new record.
type CreateRecordRequest struct {
	TTL    *int64 `json:"ttl,omitempty"`
	ZoneID string `json:"zone_id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Value  string `json:"value"`
}

// RecordsResponse represents a response from the API containing a list of records.
type RecordsResponse struct {
	Records []Record `json:"records"`
}

// RecordResponse represents a response from the API containing only one record.
type RecordResponse struct {
	Record Record `json:"record"`
}

// GetRecordByName reads the current state of a DNS Record with a given name and zone id.
func (c *Client) GetRecordByName(ctx context.Context, zoneID, name string) (*Record, error) {
	resp, err := c.request(ctx, http.MethodGet, "/api/v1/records?zone_id="+zoneID, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting records in zone %s: %w", zoneID, err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = fmt.Errorf("error closing response body: %w", cerr)
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		var response *RecordsResponse

		err = readAndParseJSONBody(resp, &response)
		if err != nil {
			return nil, err
		}

		if len(response.Records) == 0 {
			return nil, fmt.Errorf("it seems there are no records in zone %s at all", zoneID)
		}

		for _, record := range response.Records {
			if record.Name == name {
				return &record, nil
			}
		}

		return nil, fmt.Errorf("there are records in zone %s, but %s isn't included", zoneID, name)
	default:
		return nil, fmt.Errorf("http status %d unhandled", resp.StatusCode)
	}
}

// GetRecordsByZoneID reads all records in a given zone.
func (c *Client) GetRecordsByZoneID(ctx context.Context, zoneID string) (*[]Record, error) {
	resp, err := c.request(ctx, http.MethodGet, "/api/v1/records?zone_id="+zoneID, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting records in zone %s: %w", zoneID, err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = fmt.Errorf("error closing response body: %w", cerr)
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		var response *RecordsResponse

		err = readAndParseJSONBody(resp, &response)
		if err != nil {
			return nil, err
		}

		return &response.Records, nil
	default:
		return nil, fmt.Errorf("http status %d unhandled", resp.StatusCode)
	}
}

// GetRecord reads the current state of a DNS Record.
func (c *Client) GetRecord(ctx context.Context, recordID string) (*Record, error) {
	resp, err := c.request(ctx, http.MethodGet, "/api/v1/records/"+recordID, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting record %s: %w", recordID, err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = fmt.Errorf("error closing response body: %w", cerr)
		}
	}()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, nil
	case http.StatusOK:
		var response *RecordResponse

		err = readAndParseJSONBody(resp, &response)
		if err != nil {
			return nil, fmt.Errorf("error Reading json response of get record %s request: %s", recordID, err)
		}

		return &response.Record, nil
	default:
		return nil, fmt.Errorf("http status %d unhandled", resp.StatusCode)
	}
}

// CreateRecord create a new DNS records.
func (c *Client) CreateRecord(ctx context.Context, opts CreateRecordOpts) (*Record, error) {
	reqBody := CreateRecordRequest(opts)

	resp, err := c.request(ctx, http.MethodPost, "/api/v1/records", reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating record in zone %s: %w", opts.ZoneID, err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = fmt.Errorf("error closing response body: %w", cerr)
		}
	}()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, nil
	case http.StatusOK:
		var response RecordResponse

		err = readAndParseJSONBody(resp, &response)
		if err != nil {
			return nil, err
		}

		return &response.Record, nil
	default:
		return nil, fmt.Errorf("http status %d unhandled", resp.StatusCode)
	}
}

// DeleteRecord deletes a given record.
func (c *Client) DeleteRecord(ctx context.Context, id string) error {
	resp, err := c.request(ctx, http.MethodDelete, "/api/v1/records/"+id, nil)
	if err != nil {
		return fmt.Errorf("error deleting record %s: %s", id, err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
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

// UpdateRecord create a new DNS records.
func (c *Client) UpdateRecord(ctx context.Context, record Record) (*Record, error) {
	resp, err := c.request(ctx, http.MethodPut, "/api/v1/records/"+record.ID, record)
	if err != nil {
		return nil, fmt.Errorf("error updating record %s: %s", record.ID, err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = fmt.Errorf("error closing response body: %w", cerr)
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		var response RecordResponse

		err = readAndParseJSONBody(resp, &response)
		if err != nil {
			return nil, err
		}

		return &response.Record, nil
	default:
		return nil, fmt.Errorf("http status %d unhandled", resp.StatusCode)
	}
}
