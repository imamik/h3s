package cli

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	// Test basic command creation
	cmd := NewCommand(CommandConfig{
		Use:   "test",
		Short: "Test command",
		Long:  "Test command long description",
		RunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	})

	assert.Equal(t, "test", cmd.Use)
	assert.Equal(t, "Test command", cmd.Short)
	assert.Equal(t, "Test command long description", cmd.Long)
	assert.NotNil(t, cmd.RunE)

	// Test with subcommands
	subCmd := &cobra.Command{
		Use:   "sub",
		Short: "Subcommand",
	}
	cmd = NewCommand(CommandConfig{
		Use:         "test",
		Short:       "Test command",
		Subcommands: []*cobra.Command{subCmd},
	})

	assert.Len(t, cmd.Commands(), 1)
	assert.Equal(t, "sub", cmd.Commands()[0].Use)

	// Test with flags
	cmd = NewCommand(CommandConfig{
		Use:   "test",
		Short: "Test command",
		Flags: []Flag{
			{
				Name:         "string-flag",
				Shorthand:    "s",
				Value:        "",
				DefaultValue: "default",
				Usage:        "String flag",
			},
			{
				Name:         "bool-flag",
				Shorthand:    "b",
				Value:        false,
				DefaultValue: true,
				Usage:        "Bool flag",
			},
			{
				Name:         "int-flag",
				Shorthand:    "i",
				Value:        0,
				DefaultValue: 42,
				Usage:        "Int flag",
			},
			{
				Name:         "slice-flag",
				Shorthand:    "l",
				Value:        []string{},
				DefaultValue: []string{"one", "two"},
				Usage:        "Slice flag",
			},
		},
	})

	// Check string flag
	flag := cmd.Flags().Lookup("string-flag")
	assert.NotNil(t, flag)
	assert.Equal(t, "default", flag.DefValue)

	// Check bool flag
	flag = cmd.Flags().Lookup("bool-flag")
	assert.NotNil(t, flag)
	assert.Equal(t, "true", flag.DefValue)

	// Check int flag
	flag = cmd.Flags().Lookup("int-flag")
	assert.NotNil(t, flag)
	assert.Equal(t, "42", flag.DefValue)

	// Check slice flag
	flag = cmd.Flags().Lookup("slice-flag")
	assert.NotNil(t, flag)
	assert.Equal(t, "[one,two]", flag.DefValue)
}
