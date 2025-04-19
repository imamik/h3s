// Package cli provides utilities for creating and managing CLI commands
package cli

import (
	"github.com/spf13/cobra"
)

// CommandConfig holds configuration for creating a command
type CommandConfig struct {
	Use         string
	Short       string
	Long        string
	RunE        func(*cobra.Command, []string) error
	Args        cobra.PositionalArgs
	Subcommands []*cobra.Command
	Flags       []Flag
	Hidden      bool
}

// Flag represents a command flag
type Flag struct {
	Name         string
	Shorthand    string
	Value        interface{}
	DefaultValue interface{}
	Usage        string
	Required     bool
}

// NewCommand creates a new cobra command with standardized configuration
func NewCommand(config CommandConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:    config.Use,
		Short:  config.Short,
		Long:   config.Long,
		RunE:   config.RunE,
		Args:   config.Args,
		Hidden: config.Hidden,
	}

	// Add subcommands
	for _, subCmd := range config.Subcommands {
		cmd.AddCommand(subCmd)
	}

	// Add flags
	for _, flag := range config.Flags {
		addFlag(cmd, flag)
	}

	return cmd
}

// addFlag adds a flag to a command based on its type
func addFlag(cmd *cobra.Command, flag Flag) {
	switch flag.Value.(type) {
	case bool:
		defaultValue, _ := flag.DefaultValue.(bool)
		cmd.Flags().BoolP(flag.Name, flag.Shorthand, defaultValue, flag.Usage)
	case int:
		defaultValue, _ := flag.DefaultValue.(int)
		cmd.Flags().IntP(flag.Name, flag.Shorthand, defaultValue, flag.Usage)
	case string:
		defaultValue, _ := flag.DefaultValue.(string)
		cmd.Flags().StringP(flag.Name, flag.Shorthand, defaultValue, flag.Usage)
	case []string:
		defaultValue, _ := flag.DefaultValue.([]string)
		cmd.Flags().StringSliceP(flag.Name, flag.Shorthand, defaultValue, flag.Usage)
	}

	if flag.Required {
		cmd.MarkFlagRequired(flag.Name)
	}
}
