package template

import (
	"bytes"
	"text/template"
)

func CompileTemplate(
	tpl string,
	templateVars map[string]interface{},
) string {
	commandTemplate := template.Must(template.New("tpl").Parse(tpl))

	var buffer bytes.Buffer
	err := commandTemplate.Execute(&buffer, templateVars)
	if err != nil {
		panic(err)
	}

	return buffer.String()
}
