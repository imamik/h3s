// Package sshkey contains the functionality for managing Hetzner cloud SSH keys
package sshkey

import (
	"h3s/internal/cluster"
	"h3s/internal/utils/file"
	"h3s/internal/utils/logger"
	"h3s/internal/utils/resource"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// createSSHKey creates a new SSH key in Hetzner Cloud
func createSSHKey(ctx *cluster.Cluster) (*hcloud.SSHKey, error) {
	if err := validateContext(ctx); err != nil {
		return nil, err
	}

	publicKey, err := file.New(ctx.Config.SSHKeyPaths.PublicKeyPath).Load().GetString()
	if err != nil {
		return nil, err
	}

	sshKeyName := getName(ctx)
	manager := resource.NewManager[*hcloud.SSHKey](ctx, logger.SSHKey, sshKeyName)

	return manager.Create(func() (*hcloud.SSHKey, error) {
		sshKey, _, err := ctx.CloudClient.SSHKey.Create(ctx.Context, hcloud.SSHKeyCreateOpts{
			Name:      sshKeyName,
			PublicKey: publicKey,
			Labels:    ctx.GetLabels(),
		})
		return sshKey, err
	})
}

// Create creates a Hetzner cloud SSH key or returns an existing one
func Create(ctx *cluster.Cluster) (*hcloud.SSHKey, error) {
	if err := validateContext(ctx); err != nil {
		return nil, err
	}

	sshKeyName := getName(ctx)
	manager := resource.NewManager[*hcloud.SSHKey](ctx, logger.SSHKey, sshKeyName)

	return manager.GetOrCreate(
		func() (*hcloud.SSHKey, error) {
			sshKey, _, err := ctx.CloudClient.SSHKey.GetByName(ctx.Context, sshKeyName)
			return sshKey, err
		},
		func() (*hcloud.SSHKey, error) {
			return createSSHKey(ctx)
		},
	)
}
