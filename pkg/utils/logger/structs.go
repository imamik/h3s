package logger

type LogResource string

const (
	Network LogResource = "Network"
	Subnet  LogResource = "Subnet"

	SSHKey LogResource = "SSH Key"

	Pool           LogResource = "Pool"
	PlacementGroup LogResource = "Placement Group"
	Server         LogResource = "Server"
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
