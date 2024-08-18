package logger

// ResourceEvent is a log event, which contains the resource, the id, the action, the status and an error
type ResourceEvent struct {
	Resource LogResource   // Resource is the resource of the event
	ID       string        // ID is the id of the resource
	Action   interface{}   // Action is the action of the event
	Status   LogCrudStatus // Status is the status of the event
	Err      []any         // Err is the error of the event
	Depth    int           // Depth is the depth of the event (0 is the root event, 1 is a child event, ...)
	IsFirst  bool          // IsFirst is true if the event is the first event in the list
	IsLast   bool          // IsLast is true if the event is the last event in the list
}

const (
	ColorGreen     = "\033[32m" // ColorGreen is the color code to print green in the terminal
	ColorRed       = "\033[31m" // ColorRed is the color code to print red in the terminal
	ColorLightGrey = "\033[37m" // ColorLightGrey is the color code to print light grey in the terminal
	ColorDefault   = ColorLightGrey
	ColorReset     = "\033[0m" // ColorReset is the color code to reset the terminal color
)

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
	Info        LogCrudStatus = "Info"
	Success     LogCrudStatus = "Success"
	Failure     LogCrudStatus = "Failure"
)
