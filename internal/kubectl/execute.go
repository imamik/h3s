package kubectl

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/file"
)

func Execute(args []string) {

	// iterate over all filteredArgs
	for i, arg := range args {
		if arg == "-f" || arg == "--filename" {
			if args[i+1][:4] == "http" {
				continue
			}
			// replace the filename with the content of the file
			content, err := file.Load(args[i+1])
			if err != nil {
				panic(err)
			}
			args[i+1] = "- <<EOF\n" + string(content) + "\nEOF"
		}
	}

	ctx := clustercontext.Context()
	err := SSH(ctx, args)
	if err != nil {
		panic(err)
	}
}
