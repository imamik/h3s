package k3s

import (
	"github.com/spf13/cobra"
)

// Flags for the releases command
var (
	prerelease bool // Include release candidates
	stable     bool // Include stable releases
	limit      int  // Limit the number of results
)

// Cmd is the main command for k3s get commands
var Cmd = &cobra.Command{
	Use:   "k3s",
	Short: "K3S utils",
}

// getK3sReleasesCmd gets available k3s releases
var getK3sReleasesCmd = &cobra.Command{
	Use:   "releases",
	Short: "Find available k3s releases",
	RunE:  runGetK3sReleases,
}

// init adds
// subcommands to get k3s command &
// adds flags to the releases command
func init() {
	// Add subcommands to the main command
	Cmd.AddCommand(getK3sReleasesCmd)

	// Add flags to the Releases command
	getK3sReleasesCmd.Flags().BoolVar(&prerelease, "rc", false, "Include release candidates")
	getK3sReleasesCmd.Flags().BoolVar(&prerelease, "prerelease", false, "Include release candidates")
	getK3sReleasesCmd.Flags().BoolVar(&stable, "stable", false, "Include stable releases")
	getK3sReleasesCmd.Flags().IntVar(&limit, "limit", 0, "Limit the number of results")
}
