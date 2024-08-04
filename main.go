package main

import (
	"fmt"
	"h3s/cmd"
	"os"
)

// main is the entry function of the application - and will initialize the command line interface
func main() {
	err := cmd.Execute()
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if err != nil {
			fmt.Println("Error: An error occurred while trying to print the error message")
		}
		os.Exit(1)
	}
}
