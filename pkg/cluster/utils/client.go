package utils

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"os"
)

func GetClient() *hcloud.Client {
	return hcloud.NewClient(hcloud.WithToken(os.Getenv("HCLOUD_TOKEN")))
}
