package install

import (
	"fmt"
	"io"
	"os/exec"

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
