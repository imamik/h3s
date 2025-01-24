package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// readAndParseJSONBody reads the response body and parses it into the given target.
func readAndParseJSONBody(resp *http.Response, target interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = fmt.Errorf("error closing response body: %w", cerr)
		}
	}()

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("error parsing JSON response: %w", err)
	}

	return nil
}
