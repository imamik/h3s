package kubectl

import (
	"fmt"
	"strings"
)

type Command struct {
	args []string
}

func NewCommand(args []string) Command {
	return Command{
		args: args,
	}
}

func (c Command) String() string {
	c.args = append([]string{"kubectl"}, c.args...)
	return strings.Join(c.args, " ")
}

// AddKubeConfigPath adds a specific kubeconfig file
func (c Command) AddKubeConfigPath(kubeConfigPath string) Command {
	// add the kubeconfig flag to the arguments
	kubeConfigStr := fmt.Sprintf(`--kubeconfig="%s"`, kubeConfigPath)
	c.args = append([]string{kubeConfigStr}, c.args...)
	return c
}

// CompileFiles resolves files provided with "-f" or "--files" into the content
func (c Command) CompileFiles() Command {
	c.args = compileFileContent(c.args)
	return c
}
