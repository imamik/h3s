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

// runKubectl is the function that is executed when the kubectl command is called - it runs kubectl commands either directly (if setup and possible) or via SSH to the first control plane server
func runKubectl(_ *cobra.Command, args []string) error {
	ctx := clustercontext.Context()

	// check if a kubeconfig file exists
	kubeConfigPath, kubeConfigExists := kubeconfig.GetPathIfExists(ctx.Config.Name)

	// Init vars for common res & err
	var res string
	var err error

	// if a kubeconfig file exists, run kubectl commands with it, otherwise run them via SSH
	if kubeConfigExists {
		cmd := buildKubeConfigCommand(kubeConfigPath, args)
		res, err = ssh2.ExecuteLocal(cmd)
	} else {
		cmd := buildSSHCommand(args)
		res, err = ssh.SSH(ctx, cmd)
	}

	if err != nil {
		return err
	}

	fmt.Println(res)
	return nil
}

// buildKubectlCommand builds a kubectl command with a specific kubeconfig file
func buildKubeConfigCommand(kubeConfigPath string, args []string) string {
	// add the kubeconfig flag to the command
	kubeConfigStr := fmt.Sprintf(`--kubeconfig="%s"`, kubeConfigPath)
	args = append([]string{"kubectl", kubeConfigStr}, args...)
	return strings.Join(args, " ")
}

// compileFileContent replaces the filename with the content of the file
func compileFileContent(args []string) []string {
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
	return args
}

func buildSSHCommand(args []string) string {
	args = compileFileContent(args)
	args = append([]string{"kubectl"}, args...)
	return strings.Join(args, " ")
}
