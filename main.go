package main

import (
	"fmt"
	"gpm/cmd"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	args := os.Args[1:]

	command := args[0]
	restArgs := args[1:]

	switch command {
	case "install":
		if len(restArgs) < 1 {
			fmt.Println("ERROR: You must provide a package name to install.")
			os.Exit(1)
		}

		packageName := restArgs[0]
		cmd.Installer(packageName)
	case "list":
		cmd.ListRepos()
	case "--help", "-h":
		cmd.Help()
	}
}
