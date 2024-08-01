package path

import "fmt"

type FileName string

const (
	CredentialFileName FileName = "credentials.yaml"
	KubeConfigFileName FileName = "kubeconfig.yaml"
)

func GetPath(projectName string, fileName FileName) string {
	return fmt.Sprintf("$HOME/.config/hcloud-k3s/%s/%s", projectName, fileName)
}
