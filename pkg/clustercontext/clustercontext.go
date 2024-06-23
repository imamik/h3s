package clustercontext

import (
	"context"
	"hcloud-k3s-cli/pkg/config"
	"hcloud-k3s-cli/pkg/config/credentials"
)

func Context() ClusterContext {
	conf := config.Load()
	projectCredentials, _ := credentials.Get(conf)

	return ClusterContext{
		Config:      conf,
		Credentials: projectCredentials,
		Client:      GetClient(projectCredentials),
		Context:     context.Background(),
		GetName:     buildGetNameFunc(conf),
		GetLabels:   buildGetLabelsFunc(conf),
	}
}
