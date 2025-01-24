// Package main is the entry point of the application - it initializes the command line interface
package main

import (
	"fmt"
	"h3s/cmd"
	"h3s/internal/version"
	"os"
)

// main is the entry function of the application - and will initialize the command line interface
func main() {
	cmd.Initialize(version.GetBuildInfo())

	if err := cmd.Cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
