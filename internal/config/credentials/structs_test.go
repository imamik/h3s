package credentials

import (
	"testing"
)

func TestProjectCredentialsRedacted(t *testing.T) {
	creds := ProjectCredentials{
		HCloudToken:     "real-hcloud-token",
		HetznerDNSToken: "real-dns-token",
		K3sToken:        "real-k3s-token",
	}
	redacted := creds.Redacted()
	if redacted.HCloudToken != redactedText || redacted.HetznerDNSToken != redactedText || redacted.K3sToken != redactedText {
		t.Errorf("Redacted credentials did not redact all fields: %+v", redacted)
	}
}
