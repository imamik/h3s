package components

import (
	"h3s/internal/utils/template"
	"strings"
)

func kubectlApply(tpl string, data map[string]interface{}) string {
	yaml := template.CompileTemplate(tpl, data)
	cmd := strings.TrimSpace(yaml)
	return "kubectl apply -f - <<EOF\n" + cmd + "\nEOF"
}
