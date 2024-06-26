package command

import (
	"bytes"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"text/template"
)

func Worker(
	ctx clustercontext.ClusterContext,
	lb *hcloud.LoadBalancer,
) string {
	var buffer bytes.Buffer

	tpldata := make(map[string]interface{})
	tpldata["K3sVersion"] = ctx.Config.K3sVersion
	tpldata["K3sToken"] = ctx.Credentials.K3sToken
	tpldata["K3sUrl"] = lb.PrivateNet[0].IP.String()

	tpl := `if lscpu | grep Vendor | grep -q Intel; then export FLANNEL_INTERFACE=ens10 ; else export FLANNEL_INTERFACE=enp7s0 ; fi && \
curl -sfL https://get.k3s.io | K3S_TOKEN="{{ .K3sToken }}" INSTALL_K3S_VERSION="{{ .K3sVersion }}" K3S_URL=https://{{ .K3sUrl }}:6443 INSTALL_K3S_EXEC="agent \
--node-name="$(hostname -f)" \
--kubelet-arg="cloud-provider=external" \
--node-ip=$(hostname -I | awk '{print $2}') \
--node-external-ip=$(hostname -I | awk '{print $1}') \
--flannel-iface=$FLANNEL_INTERFACE" sh -`

	t := template.Must(template.New("tpl").Parse(tpl))

	err := t.Execute(&buffer, tpldata)
	if err != nil {
		panic(err)
	}

	return buffer.String()
}
