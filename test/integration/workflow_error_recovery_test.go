package integration

import (
	"os"
	"os/exec"
	"testing"
)

func BenchmarkCriticalWorkflow(b *testing.B) {
	if testing.Short() || (os.Getenv("H3S_ENABLE_REAL_INTEGRATION") != "1") {
		b.Skip("Skipping critical workflow benchmark (set H3S_ENABLE_REAL_INTEGRATION=1 to enable)")
	}
	for i := 0; i < b.N; i++ {
		cmd := exec.Command("h3s", "cluster", "create", "--name=bench-cluster", "--config=testdata/cluster_bench.yaml")
		_, _ = cmd.CombinedOutput()
		cmd = exec.Command("h3s", "cluster", "delete", "--name=bench-cluster")
		_, _ = cmd.CombinedOutput()
	}
}
