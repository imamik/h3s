package kubectl

import (
	"errors"
	"strings"
	"testing"
)

func TestNewAndAddArgs(t *testing.T) {
	cmd := New("get", "pods").AddArgs("-n", "default")
	if len(cmd.args) != 4 {
		t.Errorf("expected 4 args, got %d", len(cmd.args))
	}
}

func TestString_Success(t *testing.T) {
	cmd := New("get", "pods")
	str, err := cmd.String()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(str, "kubectl get pods") {
		t.Errorf("unexpected command string: %q", str)
	}
}

func TestString_WithError(t *testing.T) {
	cmd := New("get").addErrors(errors.New("fail"))
	_, err := cmd.String()
	if err == nil || !strings.Contains(err.Error(), "fail") {
		t.Errorf("expected error containing 'fail', got: %v", err)
	}
}

func TestNamespace(t *testing.T) {
	cmd := New("get", "svc").Namespace("foo")
	if !contains(cmd.args, "-n") || !contains(cmd.args, "foo") {
		t.Error("namespace args missing")
	}
}

func TestAddKubeConfigPath(t *testing.T) {
	cmd := New("get", "pods").AddKubeConfigPath("/tmp/kubeconfig.yaml")
	found := false
	for _, arg := range cmd.args {
		if strings.Contains(arg, "--kubeconfig=") {
			found = true
		}
	}
	if !found {
		t.Error("expected --kubeconfig flag in args")
	}
}

func TestWaitForEstablished(t *testing.T) {
	cmd := New().WaitForEstablished("crd/foo", "crd/bar")
	waitArgs := []string{"wait", "--for=condition=established", "--timeout=30s", "crd/foo crd/bar"}
	for _, wa := range waitArgs {
		if !contains(cmd.args, wa) {
			t.Errorf("expected arg %q in args", wa)
		}
	}
}

func TestGetResource(t *testing.T) {
	cmd := New().GetResource("pods")
	if !contains(cmd.args, "get") || !contains(cmd.args, "pods") {
		t.Error("expected 'get' and resource in args")
	}
}

func TestDevNull(t *testing.T) {
	cmd := New().DevNull()
	if !contains(cmd.args, ">/dev/null") || !contains(cmd.args, "2>&1") {
		t.Error("expected dev null redirection in args")
	}
}

func TestApplyTemplate_Error(t *testing.T) {
	cmd := New()
	// Use invalid template to trigger error
	cmd.ApplyTemplate("{{ .NotExist }", map[string]interface{}{})
	if cmd.Error() == nil {
		t.Error("expected error from invalid template")
	}
}

func TestEmbedFileContent_Error(t *testing.T) {
	cmd := New("apply", "-f", "nonexistent.yaml")
	cmd.EmbedFileContent()
	if cmd.Error() == nil {
		t.Error("expected error for missing file")
	}
}

func TestAddArgs_Invalid(t *testing.T) {
	cmd := New().AddArgs()
	if len(cmd.args) != 0 {
		t.Error("expected no args when adding none")
	}
}

func TestString_ErrorFormatting(t *testing.T) {
	cmd := New().addErrors(errors.New("err1"), errors.New("err2"))
	_, err := cmd.String()
	if err == nil || !strings.Contains(err.Error(), "err1") || !strings.Contains(err.Error(), "err2") {
		t.Error("error formatting did not include all errors")
	}
}

func contains(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
