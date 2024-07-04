package userdata

import "hcloud-k3s-cli/internal/utils/template"

func generateWriteFilesCommon(config CloudInitConfig) string {
	return template.CompileTemplate(`
# Script to rename the private interface to eth1 and unify NetworkManager connection naming
- path: /etc/cloud/rename_interface.sh
  content: |
    #!/bin/bash
    set -euo pipefail

    sleep 11

    INTERFACE=$(ip link show | awk '/^3:/{print $2}' | sed 's/://g')
    MAC=$(cat /sys/class/net/$INTERFACE/address)

    cat <<EOF > /etc/udev/rules.d/70-persistent-net.rules
    SUBSYSTEM=="net", ACTION=="add", DRIVERS=="?*", ATTR{address}=="$MAC", NAME="eth1"
    EOF

    ip link set $INTERFACE down
    ip link set $INTERFACE name eth1
    ip link set eth1 up

    eth0_connection=$(nmcli -g GENERAL.CONNECTION device show eth0)
    nmcli connection modify "$eth0_connection" \
      con-name eth0 \
      connection.interface-name eth0

    eth1_connection=$(nmcli -g GENERAL.CONNECTION device show eth1)
    nmcli connection modify "$eth1_connection" \
      con-name eth1 \
      connection.interface-name eth1

    systemctl restart NetworkManager
  permissions: "0744"

# Disable ssh password authentication
- content: |
    Port {{.SSHPort}}
    PasswordAuthentication no
    X11Forwarding no
    MaxAuthTries {{.SSHMaxAuthTries}}
    AllowTcpForwarding no
    AllowAgentForwarding no
    AuthorizedKeysFile .ssh/authorized_keys
  path: /etc/ssh/sshd_config.d/kube-hetzner.conf

# Set reboot method as "kured"
- content: |
    REBOOT_METHOD=kured
  path: /etc/transactional-update.conf

# Create Rancher repo config
- content: |
    [rancher-k3s-common-stable]
    name=Rancher K3s Common (stable)
    baseurl=https://rpm.rancher.io/k3s/stable/common/microos/noarch
    enabled=1
    gpgcheck=1
    repo_gpgcheck=0
    gpgkey=https://rpm.rancher.io/public.key
  path: /etc/zypp/repos.d/rancher-k3s-common.repo

# Create the kube_hetzner_selinux.te file, that allows in SELinux to not interfere with various needed services
- path: /root/kube_hetzner_selinux.te
  content: |
    module kube_hetzner_selinux 1.0;

    require {
        type kernel_t, bin_t, kernel_generic_helper_t, iscsid_t, iscsid_exec_t, var_run_t, var_lib_t,
            init_t, unlabeled_t, systemd_logind_t, systemd_hostnamed_t, container_t,
            cert_t, container_var_lib_t, etc_t, usr_t, container_file_t, container_log_t,
            container_share_t, container_runtime_exec_t, container_runtime_t, var_log_t, proc_t, io_uring_t, fuse_device_t, http_port_t,
            container_var_run_t;
        class key { read view };
        class file { open read execute execute_no_trans create link lock rename write append setattr unlink getattr watch };
        class sock_file { watch write create unlink };
        class unix_dgram_socket create;
        class unix_stream_socket { connectto read write };
        class dir { add_name create getattr link lock read rename remove_name reparent rmdir setattr unlink search write watch };
        class lnk_file { read create };
        class system module_request;
        class filesystem associate;
        class bpf map_create;
        class io_uring sqpoll;
        class anon_inode { create map read write };
        class tcp_socket name_connect;
        class chr_file { open read write };
    }

    #============= kernel_generic_helper_t ==============
    allow kernel_generic_helper_t bin_t:file execute_no_trans;
    allow kernel_generic_helper_t kernel_t:key { read view };
    allow kernel_generic_helper_t self:unix_dgram_socket create;

    #============= iscsid_t ==============
    allow iscsid_t iscsid_exec_t:file execute;
    allow iscsid_t var_run_t:sock_file write;
    allow iscsid_t var_run_t:unix_stream_socket connectto;

    #============= init_t ==============
    allow init_t unlabeled_t:dir { add_name remove_name rmdir search };
    allow init_t unlabeled_t:lnk_file create;
    allow init_t container_t:file { open read };
    allow init_t container_file_t:file { execute execute_no_trans };
    allow init_t fuse_device_t:chr_file { open read write };
    allow init_t http_port_t:tcp_socket name_connect;

    #============= systemd_logind_t ==============
    allow systemd_logind_t unlabeled_t:dir search;

    #============= systemd_hostnamed_t ==============
    allow systemd_hostnamed_t unlabeled_t:dir search;

    #============= container_t ==============
    allow container_t { cert_t container_log_t }:dir read;
    allow container_t { cert_t container_log_t }:lnk_file read;
    allow container_t cert_t:file { read open };
    allow container_t container_var_lib_t:file { create open read write rename lock setattr getattr unlink };
    allow container_t etc_t:dir { add_name remove_name write create setattr watch };
    allow container_t etc_t:file { create setattr unlink write };
    allow container_t etc_t:sock_file { create unlink };
    allow container_t usr_t:dir { add_name create getattr link lock read rename remove_name reparent rmdir setattr unlink search write };
    allow container_t usr_t:file { append create execute getattr link lock read rename setattr unlink write };
    allow container_t container_file_t:file { open read write append getattr setattr lock };
    allow container_t container_file_t:sock_file watch;
    allow container_t container_log_t:file { open read write append getattr setattr watch };
    allow container_t container_share_t:dir { read write add_name remove_name };
    allow container_t container_share_t:file { read write create unlink };
    allow container_t container_runtime_exec_t:file { read execute execute_no_trans open };
    allow container_t container_runtime_t:unix_stream_socket { connectto read write };
    allow container_t kernel_t:system module_request;
    allow container_t var_log_t:dir { add_name write remove_name watch read };
    allow container_t var_log_t:file { create lock open read setattr write unlink getattr };
    allow container_t var_lib_t:dir { add_name write read };
    allow container_t var_lib_t:file { create lock open read setattr write getattr };
    allow container_t proc_t:filesystem associate;
    allow container_t self:bpf map_create;
    allow container_t self:io_uring sqpoll;
    allow container_t io_uring_t:anon_inode { create map read write };
    allow container_t container_var_run_t:dir { add_name remove_name write };
    allow container_t container_var_run_t:file { create open read rename unlink write };

{{- if .K3sRegistries}}
# Create k3s registries file
- content: {{.K3sRegistries | base64}}
  encoding: base64
  path: /etc/rancher/k3s/registries.yaml
{{- end}}

{{- if .DNSServers}}
# Set prepare for manual dns config
- content: |
    [main]
    dns=none
  path: /etc/NetworkManager/conf.d/dns.conf

- content: |
    {{- range .DNSServers}}
    nameserver {{.}}
    {{- end}}
  path: /etc/resolv.conf
  permissions: '0644'
{{- end}}
`, config)
}
