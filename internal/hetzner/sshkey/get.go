package sshkey

import (
	"errors"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
	"h3s/internal/utils/resource"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// validateContext validates that the cluster context is properly initialized
func validateContext(ctx *cluster.Cluster) error {
	if ctx == nil {
		return errors.New("sshkey: ctx is nil")
	}
	if ctx.Config == nil {
		return errors.New("sshkey: ctx.Config is nil")
	}
	if ctx.CloudClient == nil {
		return errors.New("sshkey: ctx.CloudClient is nil")
	}
	if ctx.Context == nil {
		return errors.New("sshkey: ctx.Context is nil")
	}
	return nil
}

// Get retrieves a Hetzner cloud SSH key
func Get(ctx *cluster.Cluster) (*hcloud.SSHKey, error) {
	if err := validateContext(ctx); err != nil {
		return nil, err
	}

	sshKeyName := getName(ctx)
	manager := resource.NewManager[*hcloud.SSHKey](ctx, logger.SSHKey, sshKeyName)

	return manager.Get(func() (*hcloud.SSHKey, error) {
		sshKey, _, err := ctx.CloudClient.SSHKey.GetByName(ctx.Context, sshKeyName)
		return sshKey, err
	})
}
