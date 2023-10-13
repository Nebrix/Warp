package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	version     = "0.1.0"
	colorReset  = "\033[0m"
	colorError  = "\033[31m"
	colorSucess = "\033[32m"
)

var (
	verboseFlag = false
	enableColor = true
)

func main() {
	if len(os.Args) < 2 {
		help()
		os.Exit(1)
	}

	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "-v" || arg == "--verbose" {
			verboseFlag = true
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			break
		}
	}

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "--no-color" {
			enableColor = false
			os.Args = append(os.Args[:i], os.Args...)
			break
		}
	}

	command := os.Args[1]
	switch command {
	case "install":
		if len(os.Args) < 3 {
			fmt.Println(string(colorError), "\bERROR: You must give at least one requirement to install.", string(colorReset))
			os.Exit(1)
		}
		packageName := os.Args[2]
		if verboseFlag {
			installPackageVerbose(packageName)
		} else {
			installPackage(packageName)
		}
	case "uninstall":
		if len(os.Args) < 3 {
			fmt.Println(string(colorError), "\bERROR: You must give at least one requirement to uninstall", string(colorReset))
			os.Exit(1)
		}
		packageName := os.Args[2]
		fmt.Printf("Are you sure you want to uninstall: \n%s\n(Y/n): ", packageName)
		userChoice := getUserInput()
		if userChoice == "Y" {
			uninstallPackage(packageName)
		} else if userChoice == "n" {
			return
		}
	case "search":
		listAllPackages()
	case "--help", "-h", "help":
		help()
	case "--version", "-V":
		fmt.Printf("nebrix version: %s\n", version)
	case "update":
		if len(os.Args) < 3 {
			fmt.Println(string(colorError), "\bERROR: You must give at least one requirement to update.", string(colorReset))
			os.Exit(1)
		}
		packageName := os.Args[2]
		updatePackage(packageName)
	default:
		fmt.Printf("ERROR: unknown command: %s", command)
		os.Exit(1)
	}
}

func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func uninstallPackage(packageName string) {
	fmt.Printf("Uninstalling %s\n", packageName)
	cmd := exec.Command("rm", "-rf", packageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("Successfully uninstalled: %s\n", packageName)
	}
}

func updatePackage(packageName string) {
	originalDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting original CWD: %v\n", err)
		return
	}

	err = os.Chdir(packageName)
	if err != nil {
		fmt.Printf("Error changing directory: %v\n", err)
		return
	}

	_, err = os.Getwd()
	if err != nil {
		fmt.Printf("Error getting new CWD: %v\n", err)
		return
	}

	cmdUpdate := exec.Command("git", "pull", "origin", "main")
	cmdUpdate.Stdout = os.Stdout
	cmdUpdate.Stderr = os.Stderr
	cmdUpdate.Run()
	cmdUpdate.Wait()

	err = os.Chdir(originalDir)
	if err != nil {
		fmt.Printf("Error returning to original directory: %v\n", err)
	}
}

func installPackage(packageName string) {
	cmdHTTPS := exec.Command("git", "clone", "https://github.com/Nebrix/"+packageName+".git")
	customInstallBar(packageName)
	cmdHTTPS.Stdout = nil
	cmdHTTPS.Stderr = nil
	errHTTPS := cmdHTTPS.Run()

	if errHTTPS != nil {
		fmt.Println("HTTPS Clone Error:", errHTTPS)
		fmt.Println("Trying SSH...")
		cmdSSH := exec.Command("git", "clone", "git@github.com:Nebrix/"+packageName+".git")
		cmdSSH.Stdout = nil
		cmdSSH.Stderr = nil
	}

	if enableColor {
		fmt.Println(string(colorSucess), "\bPackage", packageName, "has been successfully installed.", string(colorReset))
	} else {
		fmt.Println("Package", packageName, "has been successfully installed.")
	}
}

func installPackageVerbose(packageName string) {
	cmdHTTPS := exec.Command("git", "clone", "https://github.com/Nebrix/"+packageName+".git")
	customInstallBar(packageName)
	cmdHTTPS.Stdout = os.Stdout
	cmdHTTPS.Stderr = os.Stderr
	errHTTPS := cmdHTTPS.Run()

	if errHTTPS != nil {
		fmt.Println("HTTPS Clone Error:", errHTTPS)
		fmt.Println("Trying SSH...")
		cmdSSH := exec.Command("git", "clone", "git@github.com:Nebrix/"+packageName+".git")
		cmdSSH.Stdout = os.Stdout
		cmdSSH.Stderr = os.Stderr
	}

	if enableColor {
		fmt.Println(string(colorSucess), "\bPackage", packageName, "has been successfully installed.", string(colorReset))
	} else {
		fmt.Println("Package", packageName, "has been successfully installed.")
	}
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
  install		Install a package.
  uninstall 		Uninstall packages.
  update		Update Nebrix package.
  list	 		List Nebrix for packages.
  help       		Show help for commands.

General Options:
  -h, --help     	Show help.

  -v, --verbose		Give more output.
  -V, --version  	Show version and exit.

  --no-color		Suppress colored output.`
	fmt.Println(helpText)
}

func customInstallBar(packageName string) {
	duration := 1 * time.Second
	barWidth := 20
	sleepTime := duration / time.Duration(barWidth)

	if enableColor {
		rainbowColors := []string{"\033[31m", "\033[91m", "\033[93m", "\033[32m", "\033[34m", "\033[36m"}

		fmt.Printf("Installing %s\n", packageName)
		for i := 0; i < barWidth; i++ {
			progress := rainbowProgress(i, barWidth)
			percentage := i * 5

			fmt.Printf("\rInstalling %s %s[%s%s%s] %d%%", packageName, rainbowColors[i%len(rainbowColors)], colorReset, progress, rainbowColors[i%len(rainbowColors)], percentage)

			time.Sleep(sleepTime)
		}
		fmt.Printf("\rInstalling %s %s[%s%s%s] 100%%\n", packageName, rainbowColors[0], colorReset, rainbowProgress(barWidth, barWidth), rainbowColors[0])
	} else {
		fmt.Printf("Installing %s\n", packageName)
		for i := 0; i < barWidth; i++ {
			progress := strings.Repeat("=", i+1)
			percentage := (i + 1) * 5

			fmt.Printf("\r[%s] %d%%", progress, percentage)

			time.Sleep(sleepTime)
		}
		fmt.Printf("\r[%s] 100%%\n", strings.Repeat("=", barWidth))
	}
}

func rainbowProgress(current, total int) string {
	rainbowColors := []string{"\033[31m", "\033[91m", "\033[93m", "\033[32m", "\033[34m", "\033[36m"}
	progress := ""

	for i := 0; i < total; i++ {
		if i < current {
			progress += rainbowColors[i%len(rainbowColors)] + "="
		} else {
			progress += " "
		}
	}

	progress += colorReset
	return progress
}
