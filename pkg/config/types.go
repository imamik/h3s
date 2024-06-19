package config

type NetworkZone string

const (
	EUCentral NetworkZone = "eu-central"
	USEast    NetworkZone = "us-east"
	USWest    NetworkZone = "us-west"
)

type Location string

const (
	Nuremberg   Location = "nbg1" // DE - EU Central
	Falkenstein Location = "fsn1" // DE - EU Central
	Helsinki    Location = "hel1" // FI - EU Central
	Ashburn     Location = "ash"  // US - US East
	Hillsboro   Location = "hil"  // US - US West
)
