package kubectl

import (
	"fmt"
	"h3s/internal/utils/file"
	"h3s/internal/utils/template"
	"strings"
)

type Command struct {
	args   []string
	errors []error
}

// New creates a new kubectl command
func New(args ...string) *Command {
	return &Command{
		args: args,
	}
}

func (c *Command) AddArgs(args ...string) *Command {
	c.args = append(c.args, args...)
	return c
}

func (c *Command) addErrors(errors ...error) *Command {
	c.errors = append(c.errors, errors...)
	return c
}

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

func (c *Command) ApplyTemplate(tpl string, data map[string]interface{}) *Command {
	content, err := template.CompileTemplate(tpl, data)
	if err != nil {
		c.addErrors(err)
	}
	c.AddArgs("apply", "-f", "-", "<<EOF\n"+strings.TrimSpace(content)+"\nEOF")
	return c
}

func (c *Command) WaitForEstablished(resources ...string) *Command {
	resourceStr := strings.Join(resources, " ")
	c.AddArgs("wait", "--for=condition=established", "--timeout=30s", resourceStr)
	return c
}

func (c *Command) GetResource(resource string) *Command {
	c.AddArgs("get", resource)
	return c
}

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
