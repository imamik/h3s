package cluster

import (
	"fmt"
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/pkg/cluster"
	"hcloud-k3s-cli/pkg/config"
	"hcloud-k3s-cli/pkg/utils/yaml"
)

var Create = &cobra.Command{
	Use:   "create",
	Short: "Create a new cluster",
	Run: func(cmd *cobra.Command, args []string) {
		// Add the logic for creating a new cluster here
		fmt.Println("Creating the new cluster")

		var conf config.Config
		err := yaml.Load("hcloud-k3s.yaml", &conf)
		if err != nil {
			fmt.Println(err)
			return
		}

		cluster.Create(conf)
	},
}
