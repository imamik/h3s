// Package commands contains various commands for the microOS image
package commands

// CleanUp returns the cleanup script for the microOS image
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
