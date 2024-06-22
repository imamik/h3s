package utils

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/config"
	"hcloud-k3s-cli/pkg/config/credentials"
	"os"
)

func GetClient(conf config.Config) *hcloud.Client {
	cred, _ := credentials.Get(conf)
	err := credentials.ValidateHCloudToken(cred.HCloudToken)
	if err != nil {
		fmt.Println("Missing valid Hetzner Cloud Token")
		fmt.Println("Option 1: Use 'hcloud-k3s config credentials' command")
		fmt.Println("Option 2: Set the environment variable HCLOUD_TOKEN")
		os.Exit(1)
	}
	fmt.Printf("Using Hetzner Cloud Token: '%s...'\n", cred.HCloudToken[:10])
	return hcloud.NewClient(hcloud.WithToken(cred.HCloudToken))
}
