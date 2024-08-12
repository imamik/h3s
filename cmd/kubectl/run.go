package kubectl

import (
	"github.com/spf13/cobra"
	"h3s/internal/cluster"
	"h3s/internal/config/kubeconfig"
	"h3s/internal/utils/common"
	"h3s/internal/utils/kubectl"
	"h3s/internal/utils/ssh"
)

// runKubectl proxies kubectl commands either directly with the kubeconfig if available or via SSH to the first control plane server
func runKubectl(cmd *cobra.Command, args []string) error {
	ctx, err := cluster.Context()
	if err != nil {
		return err
	}

	// check if a kubeconfig file exists
	kubeConfigPath, kubeConfigExists := kubeconfig.GetPathIfExists()

	// Init vars for common res & err
	command := kubectl.New(args)
	var res string

	// if a kubeconfig file exists, run kubectl commands with it, otherwise run them via SSH
	if kubeConfigExists {
		command.AddKubeConfigPath(kubeConfigPath)
		res, err = ssh.ExecuteLocal(command.String())
	} else {
		err = command.CompileFiles()
		if err != nil {
			return err
		}
		res, err = common.SSH(ctx, command.String())
	}

	if err != nil {
		return err
	}

	cmd.Println(res)
	return nil
}
