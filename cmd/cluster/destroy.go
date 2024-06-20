package cluster

import (
	"fmt"
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/pkg/cluster"
	"hcloud-k3s-cli/pkg/config/load"
)

var Destroy = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy a new cluster",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Destorying the new cluster")

		conf, err := load.Load("hcloud-k3s.yaml")
		if err != nil {
			fmt.Println(err)
			return
		}

		cluster.Destroy(conf)
	},
}
