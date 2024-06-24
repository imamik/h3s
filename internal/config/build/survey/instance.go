package survey

import (
	"github.com/charmbracelet/huh"
	"hcloud-k3s-cli/internal/config"
	"log"
)

var instanceTypeOptions = []huh.Option[config.CloudInstanceType]{
	huh.NewOption("Shared vCPU Arm", config.SharedArm),
	huh.NewOption("Shared vCPU Intel", config.SharedIntel),
	huh.NewOption("Shared vCPU AMD", config.SharedAmd),
	huh.NewOption("Dedicated vCPU AMD", config.DedicatedAmd),
}

var instanceOptions = map[config.CloudInstanceType][]huh.Option[config.CloudInstance]{
	config.SharedArm: {
		huh.NewOption("CAX11	(2  vCPUs,	4  GB RAM,	40  GB SSD)", config.CAX11),
		huh.NewOption("CAX21	(4  vCPUs,	8  GB RAM,	80  GB SSD)", config.CAX21),
		huh.NewOption("CAX31	(8  vCPUs,	16 GB RAM,	160 GB SSD)", config.CAX31),
		huh.NewOption("CAX41	(16 vCPUs,	32 GB RAM,	320 GB SSD)", config.CAX41),
	},
	config.SharedIntel: {
		huh.NewOption("CX22	(2  vCPUs,	4  GB RAM,	40  GB SSD)", config.CX22),
		huh.NewOption("CX32	(4  vCPUs,	8  GB RAM,	80  GB SSD)", config.CX32),
		huh.NewOption("CX42	(8  vCPUs,	16 GB RAM,	160 GB SSD)", config.CX42),
		huh.NewOption("CX52	(16 vCPUs,	32 GB RAM,	320 GB SSD)", config.CX52),
	},
	config.SharedAmd: {
		huh.NewOption("CPX11	(2  vCPUs,	2  GB RAM,	40  GB SSD)", config.CPX11),
		huh.NewOption("CPX21	(3  vCPUs,	4  GB RAM,	80  GB SSD)", config.CPX21),
		huh.NewOption("CPX31	(4  vCPUs,	8  GB RAM,	160 GB SSD)", config.CPX31),
		huh.NewOption("CPX41	(8  vCPUs,	16 GB RAM,	240 GB SSD)", config.CPX41),
		huh.NewOption("CPX51	(16 vCPUs,	32 GB RAM,	360 GB SSD)", config.CPX51),
	},
	config.DedicatedAmd: {
		huh.NewOption("CCX13	(2  vCPUs,	8  GB RAM,	80  GB SSD)", config.CCX13),
		huh.NewOption("CCX23	(4  vCPUs,	16 GB RAM,	160 GB SSD)", config.CCX23),
		huh.NewOption("CCX33	(8  vCPUs,	32 GB RAM,	240 GB SSD)", config.CCX33),
		huh.NewOption("CCX43	(16 vCPUs,	64 GB RAM,	360 GB SSD)", config.CCX43),
		huh.NewOption("CCX53	(32 vCPUs,	128 GB RAM,	600 GB SSD)", config.CCX53),
		huh.NewOption("CCX63	(48 vCPUs,	192 GB RAM,	960 GB SSD)", config.CCX63),
	},
}

func getInstance() config.CloudInstance {
	var instance config.CloudInstance
	var instanceType config.CloudInstanceType

	err := huh.NewSelect[config.CloudInstanceType]().
		Title("Instance Type").
		Description("Architecture and type of vCPU").
		Options(instanceTypeOptions...).
		Value(&instanceType).
		Run()

	if err != nil {
		log.Println("Error selecting instance type: %v", err)
		return instance
	}

	err = huh.NewSelect[config.CloudInstance]().
		Title("Instance").
		Description("Select an instance type").
		Options(instanceOptions[instanceType]...).
		Value(&instance).
		Run()

	if err != nil {
		log.Println("Error selecting instance: %v", err)
		return instance
	}

	return instance

}
