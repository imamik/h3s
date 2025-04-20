package components

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestManifestGeneration_Valid(t *testing.T) {
	// TODO: Test that valid input produces correct YAML manifest
	t.Skip("Manifest generation logic test scaffold")
}

func TestManifestGeneration_Invalid(t *testing.T) {
	// TODO: Test that invalid/unsupported options are handled gracefully
	t.Skip("Manifest generation error handling test scaffold")
}

func TestManifestTemplates_NotEmpty(t *testing.T) {
	templates := []struct {
		name, value string
	}{
		{"CCM", Yaml.CCM},
		{"Certificate", Yaml.Certificate},
		{"CertManager", Yaml.CertManager},
		{"CertManagerHetzner", Yaml.CertManagerHetzner},
		{"CSI", Yaml.CSI},
		{"HcloudSecrets", Yaml.HcloudSecrets},
		{"K8sDashboard", Yaml.K8sDashboard},
		{"K8sDashboardConfig", Yaml.K8sDashboardConfig},
		{"K8sIngress", Yaml.K8sIngress},
		{"Traefik", Yaml.Traefik},
		{"TraefikDashboard", Yaml.TraefikDashboard},
	}
	for _, tpl := range templates {
		if tpl.value == "" {
			t.Errorf("%s template is empty", tpl.name)
		}
	}
}

func TestManifestTemplates_YamlValid(t *testing.T) {
	// Only test a couple as a sample for validity
	toTest := []struct {
		name, value string
	}{
		{"CCM", Yaml.CCM},
		{"Certificate", Yaml.Certificate},
	}
	for _, tpl := range toTest {
		var out map[string]interface{}
		err := yaml.Unmarshal([]byte(tpl.value), &out)
		if err != nil {
			t.Errorf("%s template is not valid YAML: %v", tpl.name, err)
		}
	}
}
