package clustercontext

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/config/credentials"
	"os"
)

func GetClient(creds credentials.ProjectCredentials) *hcloud.Client {
	hcloudToken := creds.HCloudToken
	err := credentials.ValidateHCloudToken(hcloudToken)
	if err != nil {
		fmt.Println("Missing valid Hetzner Cloud Token")
		fmt.Println("Option 1: Use 'hcloud-k3s config credentials' command")
		fmt.Println("Option 2: Set the environment variable HCLOUD_TOKEN")
		os.Exit(1)
	}
	fmt.Printf("Using Hetzner Cloud Token: '%s...'\n", hcloudToken[:10])
	return hcloud.NewClient(hcloud.WithToken(hcloudToken))
}