package commands

import (
	"fmt"
	"strings"
)

func Packages() string {
	neededPackages := strings.Join([]string{
		"restorecond",
		"policycoreutils",
		"policycoreutils-python-utils",
		"setools-console",
		"audit",
		"bind-utils",
		"wireguard-tools",
		"open-iscsi",
		"nfs-client",
		"xfsprogs",
		"cryptsetup",
		"lvm2",
		"git",
		"cifs-utils",
		"bash-completion",
		"mtr",
		"tcpdump",
	}, " ")

	return fmt.Sprintf(`
set -ex
echo "First reboot successful, installing needed packages..."
transactional-update --continue pkg install -y %s
transactional-update --continue shell <<- EOF
setenforce 0
rpm --import https://rpm.rancher.io/public.key
zypper install -y https://github.com/k3s-io/k3s-selinux/releases/download/v1.4.stable.1/k3s-selinux-1.4-1.sle.noarch.rpm
zypper addlock k3s-selinux
restorecon -Rv /etc/selinux/targeted/policy
restorecon -Rv /var/lib
setenforce 1
EOF
sleep 1 && udevadm settle && reboot
`, neededPackages)
}
