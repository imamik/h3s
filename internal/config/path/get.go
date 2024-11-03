package path

type FileName string

const ( // #nosec G101 -- These are configuration file names, not credentials
	SecretsFileName    FileName = "h3s-secrets.yaml"
	KubeConfigFileName FileName = "h3s-kubeconfig.yaml"
)
