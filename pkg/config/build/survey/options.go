package survey

import (
	"github.com/charmbracelet/huh"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/k3s/releases"
)

func releasesToOptions(k3sReleases []releases.Release) []huh.Option[string] {
	var options []huh.Option[string]
	for _, release := range k3sReleases {
		var releaseString = release.Name
		for i := len(releaseString); i < 16; i++ {
			releaseString += " "
		}
		releaseString += "released on " + release.FormattedDate()
		options = append(options, huh.NewOption(releaseString, release.Name))
	}
	return options
}

var networkZoneOptions = []huh.Option[hcloud.NetworkZone]{
	huh.NewOption(string(hcloud.NetworkZoneEUCentral), hcloud.NetworkZoneEUCentral),
	huh.NewOption(string(hcloud.NetworkZoneUSEast), hcloud.NetworkZoneUSEast),
	huh.NewOption(string(hcloud.NetworkZoneUSWest), hcloud.NetworkZoneUSWest),
}
