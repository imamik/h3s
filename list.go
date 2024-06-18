package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all servers",
	Run: func(cmd *cobra.Command, args []string) {
		client := hcloud.NewClient(hcloud.WithToken(os.Getenv("HCLOUD_TOKEN")))

		servers, _, err := client.Server.List(context.Background(), hcloud.ServerListOpts{})
		if err != nil {
			fmt.Printf("Error fetching servers: %v\n", err)
			return
		}

		for _, server := range servers {
			fmt.Printf("ID: %d, Name: %s, Status: %s\n", server.ID, server.Name, server.Status)
		}
	},
}
