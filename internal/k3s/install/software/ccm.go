package software

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/template"
)

func secretYaml(ctx clustercontext.ClusterContext, network *hcloud.Network) string {
	return template.CompileTemplate(`
apiVersion: "v1"
kind: "Secret"
metadata:
  namespace: 'kube-system'
  name: 'hcloud'
stringData:
  network: "{{ .NetworkName }}"
  token: "{{ .HCloudToken }}"
`, map[string]interface{}{
		"NetworkName": network.Name,
		"HCloudToken": ctx.Credentials.HCloudToken,
	})
}

const hetznerCloudControllerManagerYaml = "https://github.com/hetznercloud/hcloud-cloud-controller-manager/releases/latest/download/ccm-networks.yaml"

func settingsYaml(ctx clustercontext.ClusterContext, network *hcloud.Network) string {
	return template.CompileTemplate(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hcloud-cloud-controller-manager
  namespace: kube-system
spec:
  template:
    spec:
      containers:
        - name: hcloud-cloud-controller-manager
          command:
            - "/bin/hcloud-cloud-controller-manager"
            - "--cloud-provider=hcloud"
            - "--leader-elect=false"
            - "--allow-untagged-cloud"
            - "--allocate-node-cidrs=true"
            - "--cluster-cidr={{ .ClusterCidrIpv4 }}"
            - "--webhook-secure-port=0"
          env:
            - name: "HCLOUD_LOAD_BALANCERS_LOCATION"
              value: "{{ .LoadbalancerLocation }}"
            - name: "HCLOUD_LOAD_BALANCERS_USE_PRIVATE_IP"
              value: "true"
            - name: "HCLOUD_LOAD_BALANCERS_ENABLED"
              value: "true"
            - name: "HCLOUD_LOAD_BALANCERS_DISABLE_PRIVATE_INGRESS"
              value: "true"
`,
		map[string]interface{}{
			"LoadbalancerLocation": ctx.Config.ControlPlane.Pool.Location,
			"ClusterCidrIpv4":      network.IPRange.String(),
		})
}

func InstallHetznerCCM(
	ctx clustercontext.ClusterContext,
	network *hcloud.Network,
	proxy *hcloud.Server,
	remote *hcloud.Server,
) {
	ApplyDynamicFile(ctx, proxy, remote, secretYaml(ctx, network))
	ApplyYaml(ctx, proxy, remote, hetznerCloudControllerManagerYaml)
	ApplyDynamicFile(ctx, proxy, remote, settingsYaml(ctx, network))
}
