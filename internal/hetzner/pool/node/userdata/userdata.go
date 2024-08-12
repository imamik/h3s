package userdata

import (
	"h3s/internal/utils/template"
)

type CloudInitConfig struct {
	Hostname          string
	SSHAuthorizedKeys []string
	SwapSize          string
	SSHPort           int
	SSHMaxAuthTries   int
	K3sRegistries     string
	DNSServers        []string
}

func GenerateCloudInitConfig(config CloudInitConfig) string {
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
		"WriteFilesCommon":  generateWriteFilesCommon(config),
		"RunCmdCommon":      generateRunCmdCommon(config),
	})
}
