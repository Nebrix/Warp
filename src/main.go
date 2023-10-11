package main

import (
	"encoding/json"
	"fmt"
	"go/build"
	"io/ioutil"
	"nebrix-package/src/version"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	exitSuccess    = 0
	repositoryLink = "https://github.com/Nebrix/Nebrix-PackageManager.git"
	formatString   = `const Version = "%s"`
)

func main() {
	if len(os.Args) < 2 {
		help()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please specify a package to install.")
			os.Exit(1)
		}
		packageName := os.Args[2]
		installPackage(packageName)
	case "list":
		listAllPackages()
	case "--help", "-h", "help":
		help()
	case "--version", "-v":
		fmt.Println("nebrix version:", version.Version)
	case "update":
		update()
	default:
		fmt.Println("Error: Unknown command:", command)
		os.Exit(1)
	}
}

func update() {
	// Use the go/build package to read the version from version/version.go in the cache
	cacheVersionPkg, err := build.Import("nebrix-package/src/version", "", build.FindOnly)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Create a temporary cache directory
	cacheDir, err := os.MkdirTemp("", "cache")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Clone the repository into the cache directory
	cmd := exec.Command("git", "clone", repositoryLink, cacheDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Define paths to version.go in both the cache and the project
	cacheVersionPath := filepath.Join(cacheVersionPkg.Dir, "version.go")
	projectVersionPath := filepath.Join("nebrix-package/src/version", "version.go")

	// Read the version directly from the version.go file in the cache
	cacheVersionData, err := os.ReadFile(cacheVersionPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Read the version directly from the version.go file in the project
	projectVersionData, err := os.ReadFile(projectVersionPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Extract versions from both cache and project
	cacheVersion := extractVersion(cacheVersionData)
	projectVersion := extractVersion(projectVersionData)

	// Check if the cache version is newer
	if cacheVersion != projectVersion {
		fmt.Println("Updating to version:", cacheVersion)

		// Update the version in the project's version.go
		updatedVersionData := []byte(fmt.Sprintf("package version\n\nconst Version = \"%s\"\n", cacheVersion))
		if err := os.WriteFile(projectVersionPath, updatedVersionData, 0644); err != nil {
			fmt.Println("Error updating version in project:", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("You're already up to date!")
	}

	// Clean up the cache
	if err := os.RemoveAll(cacheDir); err != nil {
		fmt.Println("Error cleaning up cache:", err)
	}

	os.Exit(exitSuccess)
}

func extractVersion(data []byte) string {
	// Extract the version string from the data
	versionString := string(data)
	const versionPrefix = "const Version = \""
	startIndex := strings.Index(versionString, versionPrefix)
	if startIndex < 0 {
		return ""
	}
	startIndex += len(versionPrefix)
	endIndex := strings.Index(versionString[startIndex:], "\"")
	if endIndex < 0 {
		return ""
	}
	endIndex += startIndex
	return versionString[startIndex:endIndex]
}

func installPackage(packageName string) {
	cmdHTTPS := exec.Command("git", "clone", "https://github.com/Nebrix/"+packageName+".git")
	cmdHTTPS.Stdout = os.Stdout
	cmdHTTPS.Stderr = os.Stderr
	errHTTPS := cmdHTTPS.Run()

	if errHTTPS != nil {
		fmt.Println("HTTPS Clone Error:", errHTTPS)
		fmt.Println("Trying SSH...")
		cmdSSH := exec.Command("git", "clone", "git@github.com:Nebrix/"+packageName+".git")
		cmdSSH.Stdout = os.Stdout
		cmdSSH.Stderr = os.Stderr
		errSSH := cmdSSH.Run()

		if errSSH != nil {
			fmt.Println("SSH Clone Error:", errSSH)
			fmt.Println("Trying GitHub CLI...")
			cmdCLI := exec.Command("gh", "repo", "clone", "Nebrix/"+packageName)
			cmdCLI.Stdout = os.Stdout
			cmdCLI.Stderr = os.Stderr
			errCLI := cmdCLI.Run()

			if errCLI != nil {
				fmt.Println("GitHub CLI Clone Error:", errCLI)
				fmt.Println("Failed to clone package", packageName)
				return
			}
		}
	}

	fmt.Println("Package", packageName, "has been successfully installed.")
}

type Repository struct {
	Name string `json:"name"`
}

func getFileSize(username, repoName string) (int64, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", username, repoName)

	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("GitHub API returned non-OK status code: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}

	size, ok := data["size"].(float64)
	if !ok {
		return 0, fmt.Errorf("failed to parse file size from API response")
	}

	return int64(size), nil
}

func formatFileSize(size int64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
		tb = gb * 1024
		pb = tb * 1024
	)

	switch {
	case size >= pb:
		return fmt.Sprintf("%.2f PiB", float64(size)/float64(pb))
	case size >= tb:
		return fmt.Sprintf("%.2f TiB", float64(size)/float64(tb))
	case size >= gb:
		return fmt.Sprintf("%.2f GiB", float64(size)/float64(gb))
	case size >= mb:
		return fmt.Sprintf("%.2f MiB", float64(size)/float64(mb))
	case size >= kb:
		return fmt.Sprintf("%.2f KiB", float64(size)/float64(kb))
	default:
		return fmt.Sprintf("%d bytes", size)
	}
}

func listAllPackages() error {
	username := "Nebrix"
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", username)

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API returned non-OK status code: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var repositories []Repository
	if err := json.Unmarshal(body, &repositories); err != nil {
		return err
	}

	fmt.Println("Public Repositories for Nebrix")

	for _, repo := range repositories {
		fileSize, err := getFileSize(username, repo.Name)
		if err != nil {
			return err
		}
		fmt.Printf("%s -> %s\n", repo.Name, formatFileSize(fileSize))
	}

	return nil
}

func help() {
	helpText := `
Usage:
  nebrix <command> [options]
Commands:
  add        Install a package.
  list       List installed packages.
  help       Show help for commands.
Options:
  -h, --help     Show help.
  -v, --version  Show version and exit.
`
	fmt.Println(helpText)
}
