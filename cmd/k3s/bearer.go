package k3s

import (
	"fmt"
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/k3s/bearer"
	"os/exec"
)

var Bearer = &cobra.Command{
	Use:   "bearer",
	Short: "Get the bearer token for the k3s dashboard",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := clustercontext.Context()
		bearer, err := bearer.GetBearerToken(ctx, "kubernetes-dashboard", "admin-user", 24)
		if err != nil {
			fmt.Printf("Failed to get bearer token: %v\n", err)
			return
		}

		// Correctly handling the bearer token with special characters
		copyCmd := exec.Command("sh", "-c", fmt.Sprintf("printf '%%s' \"%s\" | pbcopy", bearer))
		if err := copyCmd.Run(); err != nil {
			fmt.Printf("Failed to copy bearer token to clipboard: %v\n", err)
		} else {
			fmt.Println("Bearer token copied to clipboard.")
		}
	},
}
