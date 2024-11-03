// Package template contains the functionality for compiling templates
package template

import (
	"bytes"
	"fmt"
	"h3s/internal/utils/encode"
	"strings"
	"text/template"
)

// CompileTemplate compiles a template string with the provided variables and returns the compiled template as a string.
func CompileTemplate(templateStr string, templateVars interface{}) (string, error) {
	// Define a new template with all functions and template to parse
	tpl, err := template.
		New("tpl").
		Funcs(template.FuncMap{"base64": encode.ToBase64}).
		Parse(templateStr)
	if err != nil {
		return "", err
	}
	tpl = template.Must(tpl, err)

	// Compile the template with the provided variables to a buffer
	var buffer bytes.Buffer
	if err := tpl.Execute(&buffer, templateVars); err != nil {
		return "", err
	}

	// Check if the template contains <no value>
	str := buffer.String()
	if strings.Contains(str, "<no value>") {
		return "", fmt.Errorf("template contains <no value>")
	}

	return str, nil
}
