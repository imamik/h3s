package userdata

import (
	"h3s/internal/utils/template"
)

// GenerateRunCmdCommon generates the common runcmd section of the cloud-init user data
func GenerateRunCmdCommon(config CloudInitConfig) (string, error) {
	return template.CompileTemplate(`
# ensure that /var uses full available disk size, thanks to btrfs this is easy
- [btrfs, 'filesystem', 'resize', 'max', '/var']

{{- if ne .SSHPort 22}}
# SELinux permission for the SSH alternative port.
- [semanage, port, '-a', '-t', ssh_port_t, '-p', tcp, '{{.SSHPort}}']
{{- end}}

# Create and apply the necessary SELinux module for kube-hetzner
- [checkmodule, '-M', '-m', '-o', '/root/kube_hetzner_selinux.mod', '/root/kube_hetzner_selinux.te']
- ['semodule_package', '-o', '/root/kube_hetzner_selinux.pp', '-m', '/root/kube_hetzner_selinux.mod']
- [semodule, '-i', '/root/kube_hetzner_selinux.pp']
- [setsebool, '-P', 'virt_use_samba', '1']
- [setsebool, '-P', 'domain_kernel_load_modules', '1']

# Disable rebootmgr service as we use kured instead
- [systemctl, disable, '--now', 'rebootmgr.service']

{{- if .DNSServers}}
# Set the dns manually
- [systemctl, 'reload', 'NetworkManager']
{{- end}}

# Bounds the amount of logs that can survive on the system
- [sed, '-i', 's/#SystemMaxUse=/SystemMaxUse=3G/g', /etc/systemd/journald.conf]
- [sed, '-i', 's/#MaxRetentionSec=/MaxRetentionSec=1week/g', /etc/systemd/journald.conf]

# Reduces the default number of snapshots from 2-10 number limit, to 4 and from 4-10 number limit important, to 2
- [sed, '-i', 's/NUMBER_LIMIT="2-10"/NUMBER_LIMIT="4"/g', /etc/snapper/configs/root]
- [sed, '-i', 's/NUMBER_LIMIT_IMPORTANT="4-10"/NUMBER_LIMIT_IMPORTANT="3"/g', /etc/snapper/configs/root]

# Allow network interface
- [chmod, '+x', '/etc/cloud/rename_interface.sh']

# Restart the sshd service to apply the new config
- [systemctl, 'restart', 'sshd']

# Make sure the network is up
- [systemctl, restart, NetworkManager]
- [systemctl, status, NetworkManager]
- [ip, route, add, default, via, '10.0.0.1', ||, true]

# Cleanup some logs
- [truncate, '-s', '0', '/var/log/audit/audit.log']
`, config)
}
