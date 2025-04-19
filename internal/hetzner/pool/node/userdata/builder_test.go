package userdata

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloudInitBuilder(t *testing.T) {
	config := CloudInitConfig{
		Hostname:          "test-host",
		SwapSize:          "1G",
		SSHAuthorizedKeys: []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC..."},
		SSHPort:           22,
		SSHMaxAuthTries:   3,
	}

	builder := NewCloudInitBuilder(config)
	assert.NotNil(t, builder)

	// Test generating a specific section
	err := builder.GenerateSection(SectionHostname)
	assert.NoError(t, err)
	assert.Equal(t, "test-host", builder.sections[SectionHostname])

	// Test setting a section directly
	builder.SetSection(SectionHostname, "custom-host")
	assert.Equal(t, "custom-host", builder.sections[SectionHostname])

	// Test building the complete config
	cloudConfig, err := builder.Build()
	assert.NoError(t, err)
	assert.True(t, strings.Contains(cloudConfig, "#cloud-config"))
	assert.True(t, strings.Contains(cloudConfig, "ssh_authorized_keys:"))
	assert.True(t, strings.Contains(cloudConfig, "hostname: test-host"))
	assert.True(t, strings.Contains(cloudConfig, "- [fallocate, '-l', '1G', '/var/swapfile']"))
}

func TestCustomSectionGenerator(t *testing.T) {
	config := CloudInitConfig{
		Hostname: "test-host",
	}

	builder := NewCloudInitBuilder(config)
	assert.NotNil(t, builder)

	// Register a custom section generator
	customSection := CloudInitSection("custom")
	builder.RegisterSectionGenerator(customSection, func(config CloudInitConfig) (string, error) {
		return "custom-content: " + config.Hostname, nil
	})

	// Generate the custom section
	err := builder.GenerateSection(customSection)
	assert.NoError(t, err)
	assert.Equal(t, "custom-content: test-host", builder.sections[customSection])
}

func TestGenerateAllSections(t *testing.T) {
	config := CloudInitConfig{
		Hostname:          "test-host",
		SwapSize:          "1G",
		SSHAuthorizedKeys: []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC..."},
		SSHPort:           22,
		SSHMaxAuthTries:   3,
	}

	builder := NewCloudInitBuilder(config)
	assert.NotNil(t, builder)

	// Generate all sections
	err := builder.GenerateAllSections()
	assert.NoError(t, err)

	// Check that all registered sections were generated
	for section := range builder.sectionGen {
		_, exists := builder.sections[section]
		assert.True(t, exists, "Section %s should have been generated", section)
	}
}
