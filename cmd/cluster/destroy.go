package cluster

import (
	"fmt"
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/pkg/cluster"
	"hcloud-k3s-cli/pkg/config"
	"hcloud-k3s-cli/pkg/utils/yaml"
)

var Destroy = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy a new cluster",
	Run: func(cmd *cobra.Command, args []string) {
		var conf config.Config
		err := yaml.Load("hcloud-k3s.yaml", &conf)
		if err != nil {
			fmt.Println(err)
			return
		}
		cluster.Destroy(conf)
	},
}
