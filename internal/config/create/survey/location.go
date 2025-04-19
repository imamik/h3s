package survey

import (
	"fmt"
	"h3s/internal/config"

	"github.com/charmbracelet/huh"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func getLocation(title, description string, networkZone hcloud.NetworkZone) (config.Location, error) {
	switch networkZone {
	case hcloud.NetworkZoneEUCentral:
		var location config.Location
		err := huh.NewSelect[config.Location]().
			Title(title).
			Description(description).
			Options(
				huh.NewOption("Nürnberg (nbg1)", config.Nuernberg),
				huh.NewOption("Falkenstein (fsn1)", config.Falkenstein),
				huh.NewOption("Helsinki (hel1)", config.Helsinki),
			).
			Value(&location).
			Run()
		if err != nil {
			return "", fmt.Errorf("failed to select location: %w", err)
		}
		return location, nil
	case hcloud.NetworkZoneUSEast:
		return config.Ashburn, nil
	case hcloud.NetworkZoneUSWest:
		return config.Hillsboro, nil
	default:
		return "", fmt.Errorf("invalid network zone")
	}
}
