package k3s

import (
	"fmt"
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/pkg/k3s"
	"sort"
)

var (
	prerelease bool
	stable     bool
	limit      int
)

var FindCmd = &cobra.Command{
	Use:   "k3s versions",
	Short: "Find available k3s versions",
	Run: func(cmd *cobra.Command, args []string) {
		releases, err := k3s.GetK3sReleases()
		if err != nil {
			fmt.Println("Error fetching releases:", err)
			return
		}

		// Filter releases based on flags
		var filteredReleases []k3s.Release
		for _, version := range releases {
			if stable && !version.Prerelease && !version.Draft {
				filteredReleases = append(filteredReleases, version)
			} else if prerelease && version.Prerelease && !version.Draft {
				filteredReleases = append(filteredReleases, version)
			} else if !stable && !prerelease {
				filteredReleases = append(filteredReleases, version)
			}
		}

		// Sort by published date
		sort.Slice(filteredReleases, func(i, j int) bool {
			return filteredReleases[i].PublishedAt.Before(filteredReleases[j].PublishedAt)
		})

		// Apply limit
		if limit > 0 && limit < len(filteredReleases) {
			filteredReleases = filteredReleases[:limit]
		}

		k3s.PrintReleases(filteredReleases)

	},
}

func init() {
	FindCmd.Flags().BoolVar(&prerelease, "prerelease", false, "Include release candidates")
	FindCmd.Flags().BoolVar(&stable, "stable", false, "Include stable releases")
	FindCmd.Flags().IntVar(&limit, "limit", 0, "Limit the number of results")
}
