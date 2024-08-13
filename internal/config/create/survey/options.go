package survey

import (
	"github.com/charmbracelet/huh"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/config"
	"h3s/internal/k3s"
)

func releasesToOptions(k3sReleases []k3s.Release) []huh.Option[string] {
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

var imageOptions = []huh.Option[config.Image]{
	huh.NewOption("MicroOS", config.ImageMicroOS),
	huh.NewOption("Ubuntu 24.04", config.ImageUbuntu2404),
	huh.NewOption("Ubuntu 22.04", config.ImageUbuntu2204),
	huh.NewOption("Ubuntu 20.04", config.ImageUbuntu2004),
}
