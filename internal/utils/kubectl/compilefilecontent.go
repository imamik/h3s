package kubectl

import "h3s/internal/utils/file"

// compileFileContent replaces the filename with the content of the file
func compileFileContent(args []string) []string {
	for i, arg := range args {
		if arg == "-f" || arg == "--filename" {
			if len(args) <= i+1 {
				continue
			}
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
	return args
}
