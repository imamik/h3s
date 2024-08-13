package template

import (
	"bytes"
	"encoding/base64"
	"text/template"
)

func CompileTemplate(
	tpl string,
	templateVars interface{},
) string {
	commandTemplate := template.
		Must(template.New("tpl").
			Funcs(template.FuncMap{"base64": encodeBase64}).
			Parse(tpl))

	var buffer bytes.Buffer
	err := commandTemplate.Execute(&buffer, templateVars)
	if err != nil {
		panic(err)
	}

	return buffer.String()
}

func encodeBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
