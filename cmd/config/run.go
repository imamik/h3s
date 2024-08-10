package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"h3s/internal/config/build"
	"h3s/internal/k3s/releases"
)

// runCreate is the function that is executed when the create command is called - it creates a new h3s cluster configuration
func runCreate(_ *cobra.Command, _ []string) error {
	// Get the latest 5 k3s releases
	k3sReleases, err := releases.GetFilteredReleases(false, true, 5)

	// If there was an error, print it and return
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Build the configuration
	build.Build(k3sReleases)

	return nil
}
