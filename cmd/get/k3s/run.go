package k3s

import (
	"github.com/spf13/cobra"
	"h3s/internal/k3s/releases"
)

// runGetK3sReleases gets available k3s releases
func runGetK3sReleases(_ *cobra.Command, _ []string) error {
	r, err := releases.GetFilteredReleases(prerelease, stable, limit)
	if err != nil {
		println("Error fetching releases:", err)
		return err
	}
	releases.PrintReleases(r)
	return nil
}
