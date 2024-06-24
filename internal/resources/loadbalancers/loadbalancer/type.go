package loadbalancer

type Type string

const (
	Combined     Type = "combined"
	ControlPlane Type = "control_plane"
	Worker       Type = "worker"
)
