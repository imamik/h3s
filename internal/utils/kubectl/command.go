// Package kubectl contains the util functionality for composing kubectl commands
package kubectl

import (
	"fmt"
	"h3s/internal/utils/file"
	"h3s/internal/utils/template"
	"strings"
)

// Command represents a kubectl command
type Command struct {
	args   []string // args for the kubectl command
	errors []error  // errors that occurred during the command execution
}

// New creates a new kubectl command
func New(args ...string) *Command {
	return &Command{
		args: args,
	}
}

// AddArgs adds arguments to the kubectl command
func (c *Command) AddArgs(args ...string) *Command {
	c.args = append(c.args, args...)
	return c
}

// addErrors adds errors to the kubectl command
//
//nolint:unparam // Return value is part of fluent interface pattern for method chaining consistency
func (c *Command) addErrors(errors ...error) *Command {
	c.errors = append(c.errors, errors...)
	return c
}

// Error returns all errors that occurred during the command execution - or nil if no errors occurred
func (c *Command) Error() error {
	if len(c.errors) == 0 {
		return nil
	}
	return fmt.Errorf("%s", c.errors)
}

// String returns the kubectl command as a string
func (c *Command) String() (string, error) {
	if c.errors != nil {
		return "", c.Error()
	}
	c.args = append([]string{"kubectl"}, c.args...)
	return strings.Join(c.args, " "), nil
}

// Namespace sets the namespace for the kubectl command
func (c *Command) Namespace(namespace string) *Command {
	c.AddArgs("-n", namespace)
	return c
}

// ApplyTemplate applies a template to the kubectl command - using the template and data provided
func (c *Command) ApplyTemplate(tpl string, data map[string]interface{}) *Command {
	content, err := template.CompileTemplate(tpl, data)
	if err != nil {
		c.addErrors(err)
	}
	c.AddArgs("apply", "-f", "-", "<<EOF\n"+strings.TrimSpace(content)+"\nEOF")
	return c
}

// WaitForEstablished adds a waits for the resources to be established
func (c *Command) WaitForEstablished(resources ...string) *Command {
	resourceStr := strings.Join(resources, " ")
	c.AddArgs("wait", "--for=condition=established", "--timeout=30s", resourceStr)
	return c
}

// GetResource adds a "get" command to the kubectl command
func (c *Command) GetResource(resource string) *Command {
	c.AddArgs("get", resource)
	return c
}

// DevNull redirects the output of the kubectl command to /dev/null
func (c *Command) DevNull() *Command {
	c.AddArgs(">/dev/null", "2>&1")
	return c
}

// AddKubeConfigPath adds a specific kubeconfig file
func (c *Command) AddKubeConfigPath(kubeConfigPath string) *Command {
	kubeConfigStr := fmt.Sprintf(`--kubeconfig="%s"`, kubeConfigPath)
	c.AddArgs(kubeConfigStr)
	return c
}

// EmbedFileContent resolves files provided with "-f" or "--files" into the content
func (c *Command) EmbedFileContent() *Command {
	for i, arg := range c.args {
		if arg == "-f" || arg == "--filename" {
			if len(c.args) <= i+1 {
				continue
			}
			if c.args[i+1][:4] == "http" {
				continue
			}
			// replace the filename with the content of the file
			content, err := file.New(c.args[i+1]).Load().GetString()
			if err != nil {
				c.addErrors(err)
			}
			c.args[i+1] = "- <<EOF\n" + content + "\nEOF"
		}
	}
	return c
}
