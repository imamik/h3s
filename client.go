package main

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"os"
)

var hcloudClient *hcloud.Client

func initClient() {
	hcloudClient = hcloud.NewClient(hcloud.WithToken(os.Getenv("HCLOUD_TOKEN")))
}
