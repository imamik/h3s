package userdata

import (
	"h3s/internal/utils/template"
)

// CloudInitConfig is the configuration for the cloud-init user data
type CloudInitConfig struct {
	Hostname          string
	SwapSize          string
	K3sRegistries     string
	SSHAuthorizedKeys []string
	DNSServers        []string
	SSHPort           int
	SSHMaxAuthTries   int
}

// GenerateCloudInitConfig generates the cloud-init user data
func GenerateCloudInitConfig(config CloudInitConfig) (string, error) {
	writeFilesCommon, err := GenerateWriteFilesCommon(config)
	if err != nil {
		return "", err
	}

	runCmdCommon, err := GenerateRunCmdCommon(config)
	if err != nil {
		return "", err
	}

	return template.CompileTemplate(`#cloud-config

debug: True

write_files:

{{.WriteFilesCommon}}

# Add ssh authorized keys
ssh_authorized_keys:
{{- range .SSHAuthorizedKeys}}
  - {{.}}
{{- end}}

# Resize /var, not /, as that's the last partition in MicroOS image.
growpart:
    devices: ["/var"]

# Make sure the hostname is set correctly
hostname: {{.Hostname}}
preserve_hostname: true

runcmd:

{{.RunCmdCommon}}

{{- if .SwapSize}}
- [fallocate, '-l', '{{.SwapSize}}', '/var/swapfile']
- [chmod, '600', '/var/swapfile']
- [mkswap, '/var/swapfile']
- [swapon, '/var/swapfile']
- ["sh", "-c", "echo '/var/swapfile swap swap defaults 0 0' >> /etc/fstab"]
{{- end}}
`, map[string]interface{}{
		"Hostname":          config.Hostname,
		"SSHAuthorizedKeys": config.SSHAuthorizedKeys,
		"SwapSize":          config.SwapSize,
		"SSHPort":           config.SSHPort,
		"SSHMaxAuthTries":   config.SSHMaxAuthTries,
		"K3sRegistries":     config.K3sRegistries,
		"DNSServers":        config.DNSServers,
		"WriteFilesCommon":  writeFilesCommon,
		"RunCmdCommon":      runCmdCommon,
	})
}
