package kubectl

import (
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/ssh"
	"hcloud-k3s-cli/internal/utils/file"
	"strings"
)

var Kubectl = &cobra.Command{
	Use:                "kubectl",
	Short:              "Proxy kubectl commands via ssh to the Kubernetes API server",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {

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

		ctx := clustercontext.Context()
		ssh.SSH(ctx, strings.Join(args, " "))
	},
}
