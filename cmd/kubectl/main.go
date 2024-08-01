package kubectl

import (
	"fmt"
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config/kubeconfig"
	"hcloud-k3s-cli/internal/ssh"
	"hcloud-k3s-cli/internal/utils/file"
	ssh2 "hcloud-k3s-cli/internal/utils/ssh"
	"strings"
)

var Kubectl = &cobra.Command{
	Use:                "kubectl",
	Short:              "Proxy kubectl commands via ssh to the Kubernetes API server",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {

		ctx := clustercontext.Context()

		kubeConfigPath, kubeConfigExists := kubeconfig.GetPathIfExists(ctx.Config.Name)
		if kubeConfigExists {
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

		// iterate over all filteredArgs
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
	},
}
