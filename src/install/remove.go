package install

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getContainerID(imageName string) (string, error) {
	cmdPS := exec.Command("docker", "ps", "-a", "--quiet")
	containerIDs, err := cmdPS.Output()
	if err != nil {
		return "", err
	}

	containers := strings.Split(strings.TrimSpace(string(containerIDs)), "\n")
	for _, containerID := range containers {
		cmdInspect := exec.Command("docker", "inspect", "--format", "{{.Config.Image}}", containerID)
		containerImage, err := cmdInspect.Output()
		if err != nil {
			continue
		}
		if strings.TrimSpace(string(containerImage)) == imageName {
			return containerID, nil
		}
	}

	return "", fmt.Errorf("container with image name '%s' not found", imageName)
}

func RemovePackageGithub(packageName string) {
	cmdRemove := exec.Command("rm", "-rf", packageName)
	cmdRemove.Stdout = nil
	cmdRemove.Stderr = nil
	cmdErr := cmdRemove.Run()

	if cmdErr != nil {
		return
	}
}

func RemovePackageDocker(packageName string) {
	containerID, err := getContainerID("nebrix/" + packageName)
	if err != nil {
		fmt.Printf("Error getting container ID: %v\n", err)
		return
	}

	cmdStop := exec.Command("docker", "stop", containerID)
	cmdStop.Stdout = os.Stdout
	cmdStop.Stderr = os.Stderr
	if err := cmdStop.Run(); err != nil {
		fmt.Printf("Error stopping container: %v\n", err)
	}

	cmdRemove := exec.Command("docker", "rm", containerID)
	cmdRemove.Stdout = os.Stdout
	cmdRemove.Stderr = os.Stderr
	if err := cmdRemove.Run(); err != nil {
		fmt.Printf("Error removing container: %v\n", err)
	}

	cmdRmi := exec.Command("docker", "rmi", "nebrix/"+packageName)
	cmdRmi.Stdout = os.Stdout
	cmdRmi.Stderr = os.Stderr
	if err := cmdRmi.Run(); err != nil {
		fmt.Printf("Error removing image: %v\n", err)
	}
}
