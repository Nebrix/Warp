package install

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func DockerInstaller(packageName string) {
	cmdDOCKER := exec.Command("docker", "pull", "nebrix/"+packageName)

	cmdOut, err := cmdDOCKER.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return
	}

	errDOCKER := cmdDOCKER.Start()
	if errDOCKER != nil {
		fmt.Println("Error starting Docker pull:", errDOCKER)
		return
	}

	bar := progressbar.DefaultBytes(
		-1,
		"Pulling Docker image",
	)

	writer := io.MultiWriter(bar)

	go io.Copy(writer, cmdOut)

	errDOCKER = cmdDOCKER.Wait()
	if errDOCKER != nil {
		fmt.Println("Error:", errDOCKER)
	}

	bar.Finish()
}

func GithubInstallerHTTP(packageName string) {
	cmdHTTPS := exec.Command("git", "clone", "https://github.com/Nebrix/"+packageName+".git")

	cmdOut, err := cmdHTTPS.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return
	}

	errDOCKER := cmdHTTPS.Start()
	if errDOCKER != nil {
		fmt.Println("Error starting Docker pull:", errDOCKER)
		return
	}

	bar := progressbar.DefaultBytes(
		-1,
		"Cloning HTTP package",
	)

	writer := io.MultiWriter(bar)

	go io.Copy(writer, cmdOut)

	errDOCKER = cmdHTTPS.Wait()
	if errDOCKER != nil {
		fmt.Println("Error:", errDOCKER)
	}

	bar.Finish()
}

func GithubInstallerSSH(packageName string) {
	cmdSSH := exec.Command("git", "clone", "git@github.com:Nebrix/"+packageName+".git")
	cmdOut, err := cmdSSH.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return
	}

	errDOCKER := cmdSSH.Start()
	if errDOCKER != nil {
		fmt.Println("Error starting Docker pull:", errDOCKER)
		return
	}

	bar := progressbar.DefaultBytes(
		-1,
		"Cloning SSH package",
	)

	writer := io.MultiWriter(bar)

	go io.Copy(writer, cmdOut)

	errDOCKER = cmdSSH.Wait()
	if errDOCKER != nil {
		fmt.Println("Error:", errDOCKER)
	}

	bar.Finish()
}

func getTag(packageName string) string {
	osName := runtime.GOOS

	var cmd *exec.Cmd
	switch osName {
	case "windows":
		cmd = exec.Command("powershell", "-Command", fmt.Sprintf("Invoke-RestMethod -Uri 'https://api.github.com/repos/Nebrix/%s/releases/latest' | Select-Object -ExpandProperty tag_name", packageName))
	case "linux", "darwin":
		cmd = exec.Command("sh", "-c", `curl -s "https://api.github.com/repos/Nebrix/%s/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'`, packageName)
	default:
		fmt.Println("Unsupported operating system:", osName)
		os.Exit(1)
	}

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running git ls-remote:", err)
		os.Exit(1)
	}

	tagNumber := strings.TrimSpace(string(output))
	return tagNumber
}

func DefaultInstaller(packageName string) {
	osName := runtime.GOOS
	tag := getTag(packageName)

	switch osName {
	case "windows":
		url := fmt.Sprintf("https://github.com/Nebrix/%s/releases/download/%s/%s-%s.exe", packageName, tag, packageName, osName)
		outputFileName := fmt.Sprintf("%s-%s.exe", packageName, osName)

		cmd := exec.Command("powershell", "-Command", "Invoke-WebRequest", url, "-OutFile", outputFileName)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			fmt.Println("Error running command:", err)
			return
		}
	case "linux":
		cmd := exec.Command("curl", "-OL", "https://github.com/Nebrix/"+packageName+"/releases/download/"+tag+"/"+packageName+"-"+osName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error running curl:", err)
			return
		}
	case "darwin":
		cmd := exec.Command("curl", "-OL", "https://github.com/Nebrix/"+packageName+"/releases/download/"+tag+"/"+packageName+"-"+osName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error running curl:", err)
			return
		}
	default:
		fmt.Println("Unsupported OS")
	}
}
