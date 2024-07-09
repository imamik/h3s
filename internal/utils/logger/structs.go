package logger

type LogResource string

const (
	Cluster LogResource = "Cluster"

	Image LogResource = "Image"

	Network LogResource = "Network"
	Subnet  LogResource = "Subnet"

	SSHKey LogResource = "SSH Key"

	Pool           LogResource = "Pool"
	PlacementGroup LogResource = "Placement Group"
	Server         LogResource = "Server"
	LoadBalancer   LogResource = "Load Balancer"

	DNSZone   LogResource = "DNS Zone"
	DNSRecord LogResource = "DNS Record"
)

type LogCrudMethod string

const (
	Create LogCrudMethod = "Create"
	Get    LogCrudMethod = "Get"
	Update LogCrudMethod = "Update"
	Delete LogCrudMethod = "Delete"
)

type LogCrudStatus string

const (
	Initialized LogCrudStatus = "Init"
	Success     LogCrudStatus = "Success"
	Failure     LogCrudStatus = "Failure"
)
