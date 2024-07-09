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
  replicas: 1
  selector:
    matchLabels:
      app: hcloud-cloud-controller-manager
  template:
    metadata:
      labels:
        app: hcloud-cloud-controller-manager
    spec:
      serviceAccountName: hcloud-cloud-controller-manager
      dnsPolicy: Default
      tolerations:
        # Allow HCCM itself to schedule on nodes that have not yet been initialized by HCCM.
        - key: "node.cloudprovider.kubernetes.io/uninitialized"
          value: "true"
          effect: "NoSchedule"
        - key: "CriticalAddonsOnly"
          operator: "Exists"

        # Allow HCCM to schedule on control plane nodes.
        - key: "node-role.kubernetes.io/master"
          effect: NoSchedule
          operator: Exists
        - key: "node-role.kubernetes.io/control-plane"
          effect: NoSchedule
          operator: Exists

        - key: "node.kubernetes.io/not-ready"
          effect: "NoExecute"
      hostNetwork: true
      containers:
        - name: hcloud-cloud-controller-manager
          image: docker.io/hetznercloud/hcloud-cloud-controller-manager:v1.20.0
          args:
            - "--allow-untagged-cloud"
            - "--cloud-provider=hcloud"
            - "--route-reconciliation-period=30s"
            - "--webhook-secure-port=0"
            - "--allocate-node-cidrs=true"
            - "--cluster-cidr={{ .ClusterCidrIpv4 }}"
            - "--leader-elect=false"
          env:
            - name: "HCLOUD_LOAD_BALANCERS_LOCATION"
              value: "{{ .LoadbalancerLocation }}"
            - name: "HCLOUD_LOAD_BALANCERS_USE_PRIVATE_IP"
              value: "true"
            - name: "HCLOUD_LOAD_BALANCERS_ENABLED"
              value: "true"
            - name: "HCLOUD_LOAD_BALANCERS_DISABLE_PRIVATE_INGRESS"
              value: "true"
            - name: "HCLOUD_TOKEN"
              valueFrom:
                secretKeyRef:
                  key: token
                  name: hcloud
            - name: "HCLOUD_NETWORK"
              valueFrom:
                secretKeyRef:
                  key: network
                  name: hcloud
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
