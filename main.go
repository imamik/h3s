// Package main is the entry point of the application - it initializes the command line interface
package main

import (
	"h3s/cmd"
	"h3s/internal/version"
)

// main is the entry function of the application - and will initialize the command line interface
func main() {
	cmd.Initialize(version.GetBuildInfo())
}
