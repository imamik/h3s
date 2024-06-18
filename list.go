package main

import (
	"context"
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all servers",
	Run: func(cmd *cobra.Command, args []string) {
		servers, _, err := hcloudClient.Server.List(context.Background(), hcloud.ServerListOpts{})
		if err != nil {
			fmt.Printf("Error fetching servers: %v\n", err)
			return
		}

		if len(servers) == 0 {
			fmt.Println("No servers found")
			return
		}

		for _, server := range servers {
			fmt.Printf("ID: %d, Name: %s, Status: %s\n", server.ID, server.Name, server.Status)
		}
	},
}
