package clustercontext

import (
	"context"
	"hcloud-k3s-cli/pkg/config"
)

func Context(conf config.Config) ClusterContext {
	return ClusterContext{
		Config:    conf,
		Client:    GetClient(conf),
		Context:   context.Background(),
		GetLabels: buildGetLabelsFunc(conf),
		GetName:   buildGetNameFunc(conf),
	}
}
