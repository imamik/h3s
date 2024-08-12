package commands

func CleanUp() string {
	return `
set -ex
echo "Second reboot successful, cleaning-up..."
rm -rf /etc/ssh/ssh_host_*
echo "Make sure to use NetworkManager"
touch /etc/NetworkManager/NetworkManager.conf
sleep 1 && udevadm settle
`
}
