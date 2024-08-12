package k3s

import (
	"github.com/spf13/cobra"
	"h3s/internal/k3s"
)

// runGetK3sReleases gets available k3s releases
func runGetK3sReleases(cmd *cobra.Command, _ []string) error {
	if r, err := k3s.GetFilteredReleases(prerelease, stable, limit); err != nil {
		return err
	} else {
		k3s.PrintReleases(r)
	}
	return nil
}
