package kubectl

import (
	"github.com/spf13/cobra"
	"h3s/internal/cluster"
	"h3s/internal/config/path"
	"h3s/internal/utils/common"
	"h3s/internal/utils/file"
	"h3s/internal/utils/kubectl"
	"h3s/internal/utils/ssh"
)

// runKubectl proxies kubectl commands either directly with the kubeconfig if available or via SSH to the first control plane server
func runKubectl(cmd *cobra.Command, args []string) error {
	ctx, err := cluster.Context()
	if err != nil {
		return err
	}

	// if a kubeconfig file exists, run kubectl commands with it, otherwise run them via SSH
	kubeConfigPath, kubeConfigExists := getLocalPathIfExists()
	var res string
	if kubeConfigExists {
		res, err = runWithKubeConfig(kubeConfigPath, args)
	} else {
		res, err = runWithSSH(ctx, args)
	}

	if err != nil {
		return err
	}

	cmd.Println(res)
	return nil
}

func runWithKubeConfig(kubeConfigPath string, args []string) (string, error) {
	command := kubectl.New(args)
	command.AddKubeConfigPath(kubeConfigPath)
	return ssh.ExecuteLocal(command.String())
}

func runWithSSH(ctx *cluster.Cluster, args []string) (string, error) {
	command := kubectl.New(args)
	err := command.CompileFiles()
	if err != nil {
		return "", err
	}
	return common.SSH(ctx, command.String())
}

func getLocalPathIfExists() (string, bool) {
	p := string(path.KubeConfigFileName)
	f := file.New(p)
	absPath, _ := f.Path()
	return absPath, f.Exists()
}
