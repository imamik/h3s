package sshkey

import (
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
	"h3s/internal/utils/resource"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Delete deletes a Hetzner cloud SSH key
func Delete(ctx *cluster.Cluster) error {
	if err := validateContext(ctx); err != nil {
		return err
	}

	sshKeyName := getName(ctx)
	manager := resource.NewManager[*hcloud.SSHKey](ctx, logger.SSHKey, sshKeyName)

	// First check if the SSH key exists
	sshKey, err := Get(ctx)
	if err != nil {
		// If the key doesn't exist, consider the deletion successful
		if err.Error() == "resource is nil" {
			return nil
		}
		return err
	}

	// Delete the SSH key
	return manager.Delete(func() error {
		_, err := ctx.CloudClient.SSHKey.Delete(ctx.Context, sshKey)
		return err
	})
}
