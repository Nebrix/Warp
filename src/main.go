package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	exitSuccess    = 0
	repositoryLink = "https://github.com/Nebrix/Nebrix-PackageManager.git"
)

const (
	version = "0.0.6"
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
	case "update":
		update()
	case "--help", "-h", "help":
		help()
	case "--version", "-v":
		fmt.Println("nebrix version:", version)
	default:
		fmt.Println("Error: Unknown command:", command)
		os.Exit(1)
	}
}

func update() {
	baseDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	cacheDir := filepath.Join(baseDir, "cache")
	if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	cmd := exec.Command("git", "clone", repositoryLink, cacheDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	updateCache := filepath.Join(baseDir, "cache")

	cmd = exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = updateCache

	versionData, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	version := strings.TrimSpace(string(versionData))

	fmt.Println("Latest commit hash:", version)

	cacheFile := filepath.Join(cacheDir, "version")

	if _, err := os.Stat(cacheFile); err == nil {
		cacheData, err := os.ReadFile(cacheFile)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		if version == string(cacheData) {
			fmt.Println("You're already up to date!")
			if err := os.RemoveAll(cacheDir); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		}
	}

	if err := os.RemoveAll(cacheDir); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	cmd = exec.Command("git", "pull", repositoryLink)
	cmd.Dir = updateCache
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	os.Exit(exitSuccess)
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
  update	 Updates to current version.

Options:
  -h, --help     Show help.
  -v, --version  Show version and exit.
`
	fmt.Println(helpText)
}
