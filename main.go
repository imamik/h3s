// Package main is the entry point of the application - it initializes the command line interface
package main

import (
	"fmt"
	"h3s/cmd"
)

// main is the entry function of the application - and will initialize the command line interface
func main() {
	err := cmd.Cmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
