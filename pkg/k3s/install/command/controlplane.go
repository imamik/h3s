package command

import (
	"bytes"
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/clustercontext"
	"text/template"
)

func apiServerIP(controlPlaneNodes []*hcloud.Server) string {
	return controlPlaneNodes[0].PublicNet.IPv4.IP.String()
}

func tlsSans(controlPlaneNodes []*hcloud.Server) string {
	sans := fmt.Sprintf("--tls-san=%s", apiServerIP(controlPlaneNodes))
	for _, node := range controlPlaneNodes {
		sans += fmt.Sprintf(" --tls-san=%s", node.PrivateNet[0].IP.String())
	}
	return sans
}

func ControlPlane(
	ctx clustercontext.ClusterContext,
	controlPlaneNodes []*hcloud.Server,
	node *hcloud.Server,
) string {
	var buffer bytes.Buffer

	server := ""
	if controlPlaneNodes[0].ID == node.ID {
		server = "--cluster-init"
	} else {
		server = fmt.Sprintf("--server https://%s:6443", apiServerIP(controlPlaneNodes))
	}

	templateVars := make(map[string]interface{})
	templateVars["Server"] = server
	templateVars["TLSSans"] = tlsSans(controlPlaneNodes)
	templateVars["K3sVersion"] = ctx.Config.K3sVersion
	templateVars["K3sToken"] = ctx.Credentials.K3sToken
	//templateVars["ExtraArgs"] = fmt.Sprintf(
	//	"%s %s %s %s %s %s",
	//	c.kubeAPIServerArgsList(),
	//	c.kubeSchedulerArgsList(),
	//	c.kubeControllerManagerArgsList(),
	//	c.kubeCloudControllerManagerArgsList(),
	//	c.kubeletArgsList(),
	//	c.kubeProxyArgsList(),
	//)

	if ctx.Config.ControlPlane.AsWorkerPool {
		templateVars["Taint"] = ""
	} else {
		templateVars["Taint"] = "--node-taint CriticalAddonsOnly=true:NoExecute"
	}

	tpl := `if lscpu | grep Vendor | grep -q Intel; then export FLANNEL_INTERFACE=ens10 ; else export FLANNEL_INTERFACE=enp7s0 ; fi && \
curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION="{{ .K3sVersion }}" K3S_TOKEN="{{ .K3sToken }}" INSTALL_K3S_EXEC="server \
--disable-cloud-controller \
--disable servicelb \
--disable traefik \
--disable local-storage \
--disable metrics-server \
--write-kubeconfig-mode=644 \
--node-name="$(hostname -f)" \
--cluster-cidr=10.244.0.0/16 \
--etcd-expose-metrics=true \
{{ .FlannelWireguard }} \
--kube-controller-manager-arg="bind-address=0.0.0.0" \
--kube-proxy-arg="metrics-bind-address=0.0.0.0" \
--kube-scheduler-arg="bind-address=0.0.0.0" \
{{ .Taint }} {{ .ExtraArgs }} \
--kubelet-arg="cloud-provider=external" \
--advertise-address=$(hostname -I | awk '{print $2}') \
--node-ip=$(hostname -I | awk '{print $2}') \
--node-external-ip=$(hostname -I | awk '{print $1}') \
--flannel-iface=$FLANNEL_INTERFACE \
{{ .Server }} {{ .TLSSans }}" sh -`

	t := template.Must(template.New("tpl").Parse(tpl))

	err := t.Execute(&buffer, templateVars)
	if err != nil {
		panic(err)
	}

	return buffer.String()
}
