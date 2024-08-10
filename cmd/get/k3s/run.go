package k3s

import (
	"github.com/spf13/cobra"
	"h3s/internal/k3s/releases"
)

// runGetK3sReleases gets available k3s releases
func runGetK3sReleases(cmd *cobra.Command, _ []string) error {
	if r, err := releases.GetFilteredReleases(prerelease, stable, limit); err != nil {
		return err
	} else {
		releases.PrintReleases(r)
	}
	return nil
}
