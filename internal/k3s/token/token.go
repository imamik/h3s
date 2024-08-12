package token

import (
	"fmt"
	"h3s/internal/cluster"
	"h3s/internal/utils/common"
)

func Create(ctx *cluster.Cluster, namespace string, user string, hours int) (string, error) {
	duration := fmt.Sprintf("%dh", hours)
	cmd := fmt.Sprintf("kubectl -n %s create token %s --duration=%s", namespace, user, duration)
	bearer, err := common.SSH(ctx, cmd)
	if err != nil {
		return "", err
	}
	return bearer, nil
}
