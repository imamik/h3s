// Package version implements the version command for displaying build information
package version

import (
	"fmt"
	"h3s/internal/version"
	"os"

	"github.com/spf13/cobra"
)

// Cmd represents the version command that displays detailed version information about the h3s binary
var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Show detailed version information",
	Run: func(_ *cobra.Command, _ []string) {
		info := version.GetBuildInfo()
		fmt.Printf("h3s version %s\n", info.Version)
		fmt.Printf("Commit: %s\n", info.Commit)
		fmt.Printf("Go version: %s\n", info.GoVersion)
		os.Exit(0)
	},
}
