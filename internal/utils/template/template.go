package template

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
	"text/template"
)

func CompileTemplate(
	tpl string,
	templateVars interface{},
) (string, error) {
	commandTemplate := template.
		Must(template.New("tpl").
			Funcs(template.FuncMap{"base64": encodeBase64}).
			Parse(tpl))

	var buffer bytes.Buffer
	err := commandTemplate.Execute(&buffer, templateVars)
	if err != nil {
		return "", err
	}

	str := buffer.String()

	if strings.Contains(str, "<no value>") {
		return "", fmt.Errorf("template contains <no value>")
	}

	return buffer.String(), nil
}

func encodeBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
