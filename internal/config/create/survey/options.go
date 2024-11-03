package survey

import (
	"h3s/internal/k3s"

	"github.com/charmbracelet/huh"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func releasesToOptions(k3sReleases []k3s.Release) []huh.Option[string] {
	options := make([]huh.Option[string], 0, len(k3sReleases))
	for _, release := range k3sReleases {
		releaseString := release.Name
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
