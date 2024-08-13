package kubectl

import (
	"fmt"
	"h3s/internal/utils/file"
	"strings"
)

type Command struct {
	args []string
}

func New(args []string) *Command {
	return &Command{
		args: args,
	}
}

func (c *Command) String() string {
	c.args = append([]string{"kubectl"}, c.args...)
	return strings.Join(c.args, " ")
}

// AddKubeConfigPath adds a specific kubeconfig file
func (c *Command) AddKubeConfigPath(kubeConfigPath string) *Command {
	// add the kubeconfig flag to the arguments
	kubeConfigStr := fmt.Sprintf(`--kubeconfig="%s"`, kubeConfigPath)
	c.args = append([]string{kubeConfigStr}, c.args...)
	return c
}

// CompileFiles resolves files provided with "-f" or "--files" into the content
func (c *Command) CompileFiles() error {
	args, err := compileFileContent(c.args)
	c.args = args
	return err
}

// compileFileContent replaces the filename with the content of the file
func compileFileContent(args []string) ([]string, error) {
	for i, arg := range args {
		if arg == "-f" || arg == "--filename" {
			if len(args) <= i+1 {
				continue
			}
			if args[i+1][:4] == "http" {
				continue
			}
			// replace the filename with the content of the file
			content, err := file.New(args[i+1]).Load().GetString()
			if err != nil {
				return nil, err
			}
			args[i+1] = "- <<EOF\n" + content + "\nEOF"
		}
	}
	return args, nil
}
