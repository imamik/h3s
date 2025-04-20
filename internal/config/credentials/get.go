package credentials

import (
	"fmt"
	"h3s/internal/config/path"
	"h3s/internal/utils/file"
	"h3s/internal/validation"
)

// Get returns the ProjectCredentials from the secrets file and validates the tokens.
func Get() (*ProjectCredentials, error) {
	p := string(path.SecretsFileName)
	f := file.New(p)
	if !f.Exists() {
		return nil, fmt.Errorf("credentials file not found")
	}

	var credentials *ProjectCredentials
	err := f.Load().UnmarshalYamlTo(&credentials)
	if err != nil || credentials == nil {
		return nil, err
	}

	// Validate credentials struct
	validationErr := validation.ValidateStruct(credentials)
	if validationErr != nil {
		return nil, validationErr
	}

	// DO NOT LOG credentials directly. Use credentials.Redacted() if logging is needed.

	err = ValidateHCloudToken(credentials.HCloudToken)
	if err != nil {
		err = fmt.Errorf("missing valid Hetzner Cloud Token - Use 'h3s create credentials' command - %s", err)
		return nil, err
	}

	err = ValidateHetznerDNSToken(credentials.HetznerDNSToken)
	if err != nil {
		err = fmt.Errorf("missing valid Hetzner DNS Token - Use 'h3s create credentials' command - %s", err)
		return nil, err
	}

	return credentials, nil
}
