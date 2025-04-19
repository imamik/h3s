// Package userdata contains the functionality for generating the user data for Hetzner cloud servers
package userdata

import (
	"h3s/internal/utils/template"
)

// CloudInitSection represents a section of cloud-init configuration
type CloudInitSection string

const (
	// SectionWriteFiles is the write_files section of cloud-init
	SectionWriteFiles CloudInitSection = "write_files"
	// SectionRunCmd is the runcmd section of cloud-init
	SectionRunCmd CloudInitSection = "runcmd"
	// SectionSSHAuthorizedKeys is the ssh_authorized_keys section of cloud-init
	SectionSSHAuthorizedKeys CloudInitSection = "ssh_authorized_keys"
	// SectionHostname is the hostname section of cloud-init
	SectionHostname CloudInitSection = "hostname"
	// SectionSwap is the swap configuration section of cloud-init
	SectionSwap CloudInitSection = "swap"
	// SectionGrowPart is the growpart section of cloud-init
	SectionGrowPart CloudInitSection = "growpart"
)

// CloudInitBuilder builds cloud-init configuration
type CloudInitBuilder struct {
	config     CloudInitConfig
	sections   map[CloudInitSection]string
	sectionGen map[CloudInitSection]func(CloudInitConfig) (string, error)
}

// NewCloudInitBuilder creates a new cloud-init builder
func NewCloudInitBuilder(config CloudInitConfig) *CloudInitBuilder {
	builder := &CloudInitBuilder{
		config:     config,
		sections:   make(map[CloudInitSection]string),
		sectionGen: make(map[CloudInitSection]func(CloudInitConfig) (string, error)),
	}

	// Register default section generators
	builder.RegisterSectionGenerator(SectionWriteFiles, GenerateWriteFilesCommon)
	builder.RegisterSectionGenerator(SectionRunCmd, GenerateRunCmdCommon)
	builder.RegisterSectionGenerator(SectionSSHAuthorizedKeys, generateSSHAuthorizedKeys)
	builder.RegisterSectionGenerator(SectionHostname, generateHostname)
	builder.RegisterSectionGenerator(SectionSwap, generateSwap)
	builder.RegisterSectionGenerator(SectionGrowPart, generateGrowPart)

	return builder
}

// RegisterSectionGenerator registers a generator function for a section
func (b *CloudInitBuilder) RegisterSectionGenerator(section CloudInitSection, generator func(CloudInitConfig) (string, error)) {
	b.sectionGen[section] = generator
}

// SetSection sets the content of a section directly
func (b *CloudInitBuilder) SetSection(section CloudInitSection, content string) {
	b.sections[section] = content
}

// GenerateSection generates a section using its registered generator
func (b *CloudInitBuilder) GenerateSection(section CloudInitSection) error {
	if generator, ok := b.sectionGen[section]; ok {
		content, err := generator(b.config)
		if err != nil {
			return err
		}
		b.sections[section] = content
	}
	return nil
}

// GenerateAllSections generates all registered sections
func (b *CloudInitBuilder) GenerateAllSections() error {
	for section := range b.sectionGen {
		if err := b.GenerateSection(section); err != nil {
			return err
		}
	}
	return nil
}

// Build generates the complete cloud-init configuration
func (b *CloudInitBuilder) Build() (string, error) {
	// Generate all sections if they haven't been generated yet
	if err := b.GenerateAllSections(); err != nil {
		return "", err
	}

	// Compile the final template with all sections
	return template.CompileTemplate(`#cloud-config

debug: True

write_files:

{{.WriteFiles}}

# Add ssh authorized keys
ssh_authorized_keys:
{{.SSHAuthorizedKeys}}

# Resize /var, not /, as that's the last partition in MicroOS image.
growpart:
    devices: ["/var"]

# Make sure the hostname is set correctly
hostname: {{.Hostname}}
preserve_hostname: true

runcmd:

{{.RunCmd}}

{{.Swap}}
`, map[string]interface{}{
		"WriteFiles":        b.sections[SectionWriteFiles],
		"SSHAuthorizedKeys": b.sections[SectionSSHAuthorizedKeys],
		"Hostname":          b.config.Hostname,
		"RunCmd":            b.sections[SectionRunCmd],
		"Swap":              b.sections[SectionSwap],
	})
}

// Section generators

func generateSSHAuthorizedKeys(config CloudInitConfig) (string, error) {
	return template.CompileTemplate(`{{- range .SSHAuthorizedKeys}}
  - {{.}}
{{- end}}`, config)
}

func generateHostname(config CloudInitConfig) (string, error) {
	return config.Hostname, nil
}

func generateSwap(config CloudInitConfig) (string, error) {
	if config.SwapSize == "" {
		return "", nil
	}
	return template.CompileTemplate(`
- [fallocate, '-l', '{{.SwapSize}}', '/var/swapfile']
- [chmod, '600', '/var/swapfile']
- [mkswap, '/var/swapfile']
- [swapon, '/var/swapfile']
- ["sh", "-c", "echo '/var/swapfile swap swap defaults 0 0' >> /etc/fstab"]`, config)
}

func generateGrowPart(config CloudInitConfig) (string, error) {
	return `devices: ["/var"]`, nil
}
