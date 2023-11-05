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
			fmt.Println("ERROR: You must provide a package name to install.")
			os.Exit(1)
		}
		packageName := os.Args[2]
		if githubFlag {
			if len(os.Args) < 4 {
				fmt.Println("ERROR: You must provide a method for GitHub installation (--ssh or --http).")
				os.Exit(1)
			}
			method := os.Args[3]
			switch method {
			case "--ssh":
				install.GithubInstallerSSH(packageName)
			case "--http":
				install.GithubInstallerHTTP(packageName)
			default:
				fmt.Println("Error: Unsupported installation method.")
				os.Exit(1)
			}
		} else if dockerFlag {
			install.DockerInstaller(packageName)
		}
	case "search":
		if dockerFlag {
			helper.ListAllPackagesDocker()
		} else if githubFlag {
			helper.ListAllPackages()
		}
	case "--help", "-h":
		help()
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
  install     Install a package.
  search      List Nebrix for packages.
  help        Show help for commands.

General Options:
  -h, --help      Show help.

  -D, --docker    Runs the docker install.
  -G, --github    Runs the github install.
  --ssh, --http	  Runs the github install method.`
	fmt.Println(helpText)
}
