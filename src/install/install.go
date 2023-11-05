package install

import (
	"fmt"
	"os"
	"os/exec"
)

func DockerInstaller(packageName string) {
	cmdDOCKER := exec.Command("docker", "pull", "nebrix/"+packageName)
	cmdDOCKER.Stdout = os.Stdout
	cmdDOCKER.Stderr = os.Stderr
	errDOCKER := cmdDOCKER.Run()

	if errDOCKER != nil {
		fmt.Println("Error:", errDOCKER)
	}
}

func GithubInstallerHTTP(packageName string) {
	//	fmt.Println("Cloning HTTP package")
	cmdHTTPS := exec.Command("git", "clone", "https://github.com/Nebrix/"+packageName+".git")
	cmdHTTPS.Stdout = os.Stdout
	cmdHTTPS.Stderr = os.Stderr
	errHTTPS := cmdHTTPS.Run()

	if errHTTPS != nil {
		fmt.Println("HTTPS Clone Error:", errHTTPS)
		return
	}
}

func GithubInstallerSSH(packageName string) {
	//	fmt.Println("Cloning SSH package")
	cmdSSH := exec.Command("git", "clone", "git@github.com:Nebrix/"+packageName+".git")
	cmdSSH.Stdout = os.Stdout
	cmdSSH.Stderr = os.Stderr
	errSSH := cmdSSH.Run()

	if errSSH != nil {
		fmt.Println("SSH Clone Error:", errSSH)
		return
	}
}
