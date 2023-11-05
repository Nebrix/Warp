package main

import (
	"fmt"
	"os"
	"warpPackage/src/helper"
	"warpPackage/src/install"
)

var (
	githubFlag = false
	dockerFlag = false
)

func main() {
	if len(os.Args) < 2 {
		help()
		os.Exit(1)
	}

	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "-D" || arg == "--docker" {
			dockerFlag = true
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			break
		}
	}

	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "-G" || arg == "--github" {
			githubFlag = true
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			break
		}
	}

	command := os.Args[1]
	switch command {
	case "install":
		if len(os.Args) < 3 {
			fmt.Println("ERROR: You must give at least one requirement to install.")
			os.Exit(1)
		}
		packageName := os.Args[2]
		if dockerFlag {
			install.DockerInstaller(packageName)
		} else if githubFlag {
			install.GithubInstaller(packageName)
		}
	case "--help", "-h":
		help()
	case "search":
		helper.ListAllPackages()
		helper.ListAllPackagesDocker()
	default:
		fmt.Printf("ERROR: unknown command: %s", command)
		os.Exit(1)
	}
}

func help() {
	helpText := `
Usage:
  warp <command> [options]

Commands:
  install		Install a package.
  search	 	List Nebrix for packages.
  help       		Show help for commands.

General Options:
  -h, --help     	Show help.

  -D, --docker		Runs the docker install.
  -G, --github  	Runs the github install.`
	fmt.Println(helpText)
}
