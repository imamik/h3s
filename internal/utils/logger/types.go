package logger

const (
	ColorGreen     = "\033[32m"     // ColorGreen is the color code to print green in the terminal
	ColorRed       = "\033[31m"     // ColorRed is the color code to print red in the terminal
	ColorLightGrey = "\033[37m"     // ColorLightGrey is the color code to print light grey in the terminal
	ColorDefault   = ColorLightGrey // ColorDefault is the color code to print default in the terminal
	ColorReset     = "\033[0m"      // ColorReset is the color code to reset the terminal color
)

// LogResource is a resource type for the logger
type LogResource string

const (
	Cluster LogResource = "Cluster" // Cluster is a resource representing a Hetzner cloud cluster

	Image LogResource = "Image" // Image is a resource representing a Hetzner cloud image

	Network LogResource = "Network" // Network is a resource representing a Hetzner cloud network
	Subnet  LogResource = "Subnet"  // Subnet is a resource representing a Hetzner cloud subnet

	SSHKey LogResource = "SSH Key" // SSHKey is a resource representing a Hetzner cloud SSH key

	Pool           LogResource = "Pool"            // Pool is a resource representing a Hetzner cloud pool
	PlacementGroup LogResource = "Placement Group" // PlacementGroup is a resource representing a Hetzner cloud placement group
	Server         LogResource = "Server"          // Server is a resource representing a Hetzner cloud server
	LoadBalancer   LogResource = "Load Balancer"   // LoadBalancer is a resource representing a Hetzner cloud load balancer

	DNSZone   LogResource = "DNS Zone"   // DNSZone is a resource representing a Hetzner cloud DNS zone
	DNSRecord LogResource = "DNS Record" // DNSRecord is a resource representing a Hetzner cloud DNS record
)

// LogCrudMethod is a method type for the logger
type LogCrudMethod string

const (
	Create LogCrudMethod = "Create" // Create is a method representing a create action
	Get    LogCrudMethod = "Get"    // Get is a method representing a get action
	Update LogCrudMethod = "Update" // Update is a method representing an update action
	Delete LogCrudMethod = "Delete" // Delete is a method representing a delete action
)

// LogCrudStatus is a status type for the logger
type LogCrudStatus string

const (
	Initialized LogCrudStatus = "Init"    // Initialized is a status representing an initialized action
	Info        LogCrudStatus = "Info"    // Info is a status representing an info action
	Success     LogCrudStatus = "Success" // Success is a status representing a success action
	Failure     LogCrudStatus = "Failure" // Failure is a status representing a failure action
)
