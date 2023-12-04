package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	kb = 1024
	mb = kb * 1024
	gb = mb * 1024
	tb = gb * 1024
	pb = tb * 1024
)

type Repository struct {
	Name string `json:"name"`
}

func getFileSize(repoName string) (int64, error) {
	url := fmt.Sprintf("https://api.github.com/repos/Nebrix/%v", repoName)

	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("GitHub API returned non-OK status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
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

func ListRepos() error {
	url := fmt.Sprintf("https://api.github.com/users/Nebrix/repos")

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API returned non-OK status code: %v", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var repositories []Repository
	if err := json.Unmarshal(body, &repositories); err != nil {
		return err
	}

	fmt.Println("Public Repositories for Nebrix")

	for _, repo := range repositories {
		fileSize, err := getFileSize(repo.Name)
		if err != nil {
			return err
		}
		fmt.Printf("%v -> %v\n", repo.Name, formatFileSize(fileSize))
	}
	return nil
}

func formatFileSize(size int64) string {
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

func getTag(packageName string) string {
	osName := runtime.GOOS

	var cmd *exec.Cmd
	switch osName {
	case "windows":
		cmd = exec.Command("powershell", "-Command", fmt.Sprintf("(Invoke-RestMethod -Uri 'https://api.github.com/repos/Nebrix/%v/releases/latest').tag_name", packageName))
	case "linux", "darwin":
		cmd = exec.Command("sh", "-c", fmt.Sprintf(`curl -s "https://api.github.com/repos/Nebrix/%v/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'`, packageName))
	default:
		fmt.Println("Unsupported operating system:", osName)
		os.Exit(1)
	}

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running command:", err)
		os.Exit(1)
	}

	version := strings.TrimSpace(string(output))
	return version
}

func Installer(packageName string) {
	osName := runtime.GOOS
	systemTag := runtime.GOARCH
	tag := getTag(packageName)

	switch osName {
	case "windows":
		url := fmt.Sprintf("https://github.com/Nebrix/%v/releases/download/%v/%v-%v-%v.exe", packageName, tag, packageName, osName, systemTag)
		outputFileName := fmt.Sprintf("%v-%v-%v.exe", packageName, osName, systemTag)

		cmd := exec.Command("powershell", "-Command", "Invoke-WebRequest", url, "-OutFile", outputFileName)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			fmt.Println("Error running command:", err)
			return
		}
	case "linux", "darwin":
		cmd := exec.Command("wget", "https://github.com/Nebrix/"+packageName+"/releases/download/"+tag+"/"+packageName+"-"+osName+"-"+systemTag)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error running wget:", err)
			return
		}
	default:
		fmt.Println("Unsupported OS")
	}
}

func Help() {
	helpText := `
Usage:
  warp <command> [options]

Commands:
  install     Install a package.
  List        List Nebrix for packages.
  help        Show help for commands.

General Options:
  -h, --help    Show this help message.

Examples:
  warp install mypackage
  warp list
`
	fmt.Println(helpText)
}
