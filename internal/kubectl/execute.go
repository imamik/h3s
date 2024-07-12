package kubectl

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/file"
)

func Execute(args []string) {
	useSsh := false
	var filteredArgs []string

	for _, arg := range args {
		if arg == "--ssh" {
			useSsh = true
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
	}

	// iterate over all filteredArgs
	for i, arg := range filteredArgs {
		if arg == "-f" || arg == "--filename" {
			if filteredArgs[i+1][:4] == "http" {
				continue
			}
			// replace the filename with the content of the file
			content, err := file.Load(filteredArgs[i+1])
			if err != nil {
				panic(err)
			}
			filteredArgs[i+1] = "- <<EOF\n" + string(content) + "\nEOF"
		}
	}

	if useSsh {
		ctx := clustercontext.Context()
		err := SSH(ctx, filteredArgs)
		if err != nil {
			panic(err)
		}
	} else {
		Kubectl(filteredArgs)
	}
}
