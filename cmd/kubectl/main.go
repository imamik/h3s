package kubectl

import (
	"fmt"
	"github.com/spf13/cobra"
	"h3s/internal/clustercontext"
	"h3s/internal/config/kubeconfig"
	"h3s/internal/ssh"
	"h3s/internal/utils/file"
	ssh2 "h3s/internal/utils/ssh"
	"strings"
)

var Kubectl = &cobra.Command{
	Use:                "kubectl",
	Short:              "Run kubectl commands",
	Long:               `Run kubectl commands either directly (if setup and possible) or via SSH to the first control plane server`,
	DisableFlagParsing: true,
	Run:                runKubectl,
}

func runKubectl(cmd *cobra.Command, args []string) {
	ctx := clustercontext.Context()

	kubeConfigPath, kubeConfigExists := kubeconfig.GetPathIfExists(ctx.Config.Name)
	if kubeConfigExists {
		runWithKubeConfig(kubeConfigPath, args)
	} else {
		runWithSSH(ctx, args)
	}

}

func runWithKubeConfig(kubeConfigPath string, args []string) {
	kubeConfigStr := fmt.Sprintf(`--kubeconfig="%s"`, kubeConfigPath)
	args = append([]string{"kubectl", kubeConfigStr}, args...)
	cmd := strings.Join(args, " ")
	println(cmd)
	out, err := ssh2.ExecuteLocal(cmd)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(out)
	return
}

func runWithSSH(ctx clustercontext.ClusterContext, args []string) {
	for i, arg := range args {
		if arg == "-f" || arg == "--filename" {
			if len(args) <= i+1 {
				continue
			}
			if args[i+1][:4] == "http" {
				continue
			}
			// replace the filename with the content of the file
			content, err := file.Load(args[i+1])
			if err != nil {
				panic(err)
			}
			args[i+1] = "- <<EOF\n" + string(content) + "\nEOF"
		}
	}

	args = append([]string{"kubectl"}, args...)

	ssh.SSH(ctx, strings.Join(args, " "))
}
