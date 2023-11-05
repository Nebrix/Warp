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

func GithubInstaller(packageName string) {
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
	}
}
