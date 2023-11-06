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

	args := os.Args[1:]

	for i, arg := range args {
		switch arg {
		case "-D", "--docker":
			dockerFlag = true
			args = append(args[:i], args[i+1:]...)
		case "-G", "--github":
			githubFlag = true
			args = append(args[:i], args[i+1:]...)
		}
	}

	if len(args) < 1 {
		fmt.Println("ERROR: Missing command.")
		os.Exit(1)
	}

	command := args[0]
	restArgs := args[1:]

	switch command {
	case "install":
		if len(restArgs) < 1 {
			fmt.Println("ERROR: You must provide a package name to install.")
			os.Exit(1)
		}

		packageName := restArgs[0]

		if githubFlag {
			if len(restArgs) < 2 {
				fmt.Println("ERROR: You must provide an installation method (--ssh or --http).")
				os.Exit(1)
			}

			method := restArgs[1]

			switch method {
			case "--ssh":
				install.GithubInstallerSSH(packageName)
			case "--http":
				install.GithubInstallerHTTP(packageName)
			default:
				fmt.Println("ERROR: Unsupported installation method.")
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
	case "remove":
		if len(restArgs) < 1 {
			fmt.Println("ERROR: You must provide a package name to remove.")
			os.Exit(1)
		}

		packageName := restArgs[0]

		if githubFlag {
			install.RemovePackageGithub(packageName)
		} else if dockerFlag {
			install.RemovePackageDocker(packageName)
		}
	case "--help", "-h":
		help()
	default:
		fmt.Printf("ERROR: Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func help() {
	helpText := `
Usage:
  warp <command> [options]

Commands:
  install     Install a package.
  remove      Remove a package.
  search      List Nebrix for packages.
  help        Show help for commands.

General Options:
  -h, --help    Show this help message.
  -D, --docker  Install packages using Docker.
  -G, --github  Install packages from GitHub.
  --ssh, --http  Specify the installation method for GitHub packages (SSH or HTTP).

Examples:
  warp install mypackage -D
  warp remove anotherpackage -G
  warp search -D
`
	fmt.Println(helpText)
}
