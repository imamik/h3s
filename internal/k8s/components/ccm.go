package components

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
)

func CCMServiceAccount() string {
	return kubectlApply(`
apiVersion: v1
kind: ServiceAccount
metadata:
  name: hcloud-cloud-controller-manager
  namespace: kube-system
`, nil)
}

func CCMRoleBinding() string {
	return kubectlApply(`
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: "system:hcloud-cloud-controller-manager"
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: hcloud-cloud-controller-manager
    namespace: kube-system
`, nil)
}

func CCMSettings(ctx *cluster.Cluster, network *hcloud.Network) string {
	return kubectlApply(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hcloud-cloud-controller-manager
  namespace: kube-system
spec:
  replicas: 1
  revisionHistoryLimit: 2
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
      priorityClassName: system-cluster-critical
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
          ports:
            - name: metrics
              containerPort: 8233
          resources:
            requests:
              cpu: 100m
              memory: 50Mi
`,
		map[string]interface{}{
			"LoadbalancerLocation": ctx.Config.ControlPlane.Pool.Location,
			"ClusterCidrIpv4":      network.IPRange.String(),
		})
}
