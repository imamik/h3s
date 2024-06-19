package survey

import (
	"github.com/charmbracelet/huh"
	"hcloud-k3s-cli/pkg/config"
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

var networkZoneOptions = []huh.Option[config.NetworkZone]{
	huh.NewOption(string(config.EUCentral), config.EUCentral),
	huh.NewOption(string(config.USEast), config.USEast),
	huh.NewOption(string(config.USWest), config.USEast),
}
