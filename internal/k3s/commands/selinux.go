package commands

func SeLinux() string {
	return "/sbin/semodule -v -i /usr/share/selinux/packages/k3s.pp"
}
