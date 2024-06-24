package survey

import (
	"github.com/charmbracelet/huh"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/config"
	"log"
)

func getLocation(title string, description string, networkZone hcloud.NetworkZone) config.Location {
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
			log.Fatal(err)
			return ""
		}
		return location
	case hcloud.NetworkZoneUSEast:
		return config.Ashburn
	case hcloud.NetworkZoneUSWest:
		return config.Hillsboro
	default:
		log.Fatal("Invalid network zone")
		return ""
	}
}
