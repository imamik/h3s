package components

import _ "embed"

//go:embed csi.yaml
var csiYAML string

func CSIHelmChart() string {
	return kubectlApply(csiYAML, nil)
}
