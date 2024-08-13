package commands

func WriteImage() string {
	return `
set -ex
echo 'MicroOS image loaded, writing to disk... '
qemu-img convert -p -f qcow2 -O host_device $(ls -a | grep -ie '^opensuse.*microos.*qcow2$') /dev/sda
echo 'done. Rebooting...'
sleep 1 && udevadm settle && reboot
`
}
