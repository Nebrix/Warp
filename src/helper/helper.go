package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Repository struct {
	Name string `json:"name"`
}

type DockerHubResponse struct {
	Results []Repository `json:"results"`
}

func ListAllPackages() error {
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

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var repositories []Repository
	if err := json.Unmarshal(body, &repositories); err != nil {
		return err
	}

	fmt.Println("Public Repositories for Nebrix Github")

	for _, repo := range repositories {
		fileSize, err := getFileSize(username, repo.Name)
		if err != nil {
			return err
		}
		fmt.Printf("%s -> %s\n", repo.Name, formatFileSize(fileSize))
	}

	return nil
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

func ListAllPackagesDocker() {
	organization := "nebrix"
	url := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/", organization)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed to make the request:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Received non-OK HTTP status code: %d\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read the response body:", err)
		return
	}

	var response DockerHubResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Failed to unmarshal JSON:", err)
		return
	}

	fmt.Println("Public Repositories for Nebrix Docker")

	for _, repo := range response.Results {
		fmt.Println(repo.Name)
	}
}

func Help() {
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
