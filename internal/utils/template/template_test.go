package template

import (
	"strings"
	"testing"
)

func TestCompileTemplate_Success(t *testing.T) {
	tpl := "Hello, {{.Name}}!"
	vars := map[string]interface{}{"Name": "World"}
	out, err := CompileTemplate(tpl, vars)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "World") {
		t.Errorf("expected output to contain 'World', got: %q", out)
	}
}

func TestCompileTemplate_MissingVar(t *testing.T) {
	tpl := "Hello, {{.Missing}}!"
	vars := map[string]interface{}{"Name": "World"}
	_, err := CompileTemplate(tpl, vars)
	if err == nil || !strings.Contains(err.Error(), "<no value>") {
		t.Errorf("expected error about <no value>, got: %v", err)
	}
}

func TestCompileTemplate_WithFunc(t *testing.T) {
	tpl := "{{base64 .Data}}"
	vars := map[string]interface{}{"Data": "foo"}
	out, err := CompileTemplate(tpl, vars)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Zm9v") {
		t.Errorf("expected base64 output, got: %q", out)
	}
}
