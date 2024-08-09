package k3s

import (
	"github.com/spf13/cobra"
)

// Define the flags for the releases command
var (
	prerelease bool
	stable     bool
	limit      int
)

// Cmd defines the main command and subcommands
var Cmd = &cobra.Command{
	Use:   "k3s",
	Short: "K3S utils",
}

// releasesCmd defines the command to find available k3s releases
var releasesCmd = &cobra.Command{
	Use:   "releases",
	Short: "Find available k3s releases",
	Run:   runReleases,
}

// kubeConfigCmd defines the command to get the kubeconfig for the k3s cluster
var kubeConfigCmd = &cobra.Command{
	Use:   "kubeconfig",
	Short: "Get the kubeconfig for the k3s cluster",
	Run:   runKubeConfig,
}

// installCmd defines the command to install k3s on all servers in the cluster
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install k3s on all servers in the cluster",
	Run:   runInstall,
}

// installSoftwareCmd defines the command to install software on all servers in the cluster
var installSoftwareCmd = &cobra.Command{
	Use:   "software",
	Short: "Install software on all servers in the cluster",
	Run:   runInstallSoftware,
}

// tokenCmd defines the command to get the bearer token for the k3s dashboard
var tokenCmd = &cobra.Command{
	Use:   "bearer",
	Short: "Get the bearer token for the k3s dashboard",
	Run:   runToken,
}

// init adds
// subcommands to the main command,
// subcommands to the install command &
// flags to the releases command
func init() {
	// Add subcommands to the main command
	Cmd.AddCommand(releasesCmd)
	Cmd.AddCommand(installCmd)
	Cmd.AddCommand(tokenCmd)
	Cmd.AddCommand(kubeConfigCmd)

	// Add subcommands to the install command
	installCmd.AddCommand(installSoftwareCmd)

	// Add flags to the Releases command
	releasesCmd.Flags().BoolVar(&prerelease, "rc", false, "Include release candidates")
	releasesCmd.Flags().BoolVar(&prerelease, "prerelease", false, "Include release candidates")
	releasesCmd.Flags().BoolVar(&stable, "stable", false, "Include stable releases")
	releasesCmd.Flags().IntVar(&limit, "limit", 0, "Limit the number of results")
}
