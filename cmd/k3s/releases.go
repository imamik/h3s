package k3s

import (
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/pkg/k3s/releases"
)

var (
	prerelease bool
	stable     bool
	limit      int
)

var Releases = &cobra.Command{
	Use:   "releases",
	Short: "Find available k3s releases",
	Run: func(cmd *cobra.Command, args []string) {
		r, err := releases.GetFilteredReleases(prerelease, stable, limit)
		if err != nil {
			println("Error fetching releases:", err)
			return
		}
		releases.PrintReleases(r)
	},
}

func init() {
	Releases.Flags().BoolVar(&prerelease, "rc", false, "Include release candidates")
	Releases.Flags().BoolVar(&prerelease, "prerelease", false, "Include release candidates")
	Releases.Flags().BoolVar(&stable, "stable", false, "Include stable releases")
	Releases.Flags().IntVar(&limit, "limit", 0, "Limit the number of results")
}
