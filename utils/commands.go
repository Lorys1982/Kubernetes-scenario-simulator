package utils

import (
	"log"
	"main/configs"
	"os"
	"os/exec"
)

func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// selector:
//
//	false -> cluster deletion
//	true -> cluster creation
func clusterArgs(selector bool) []string {
	var command []string
	if selector {
		command = []string{"create", "cluster"}
	} else {
		command = []string{"delete", "cluster"}
	}

	if configs.GetClusterName() != "" {
		command = append(command, "--name", configs.GetClusterName())
	}

	if selector && configs.GetScheduler() != "" {
		command = append(command, "--config", configs.GetScheduler())
	}

	return command
}

func commandRun(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func commandCleanRun(cmd *exec.Cmd) error {
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func KwokctlCreate() {
	args := clusterArgs(true)

	cmd := exec.Command("kwokctl", args...)
	err := commandRun(cmd)
	if err != nil {
		KwokctlDelete()
		log.Fatal(err.Error())
	}
}

func KwokctlDelete() {
	args := clusterArgs(false)

	cmd := exec.Command("kwokctl", args...)
	err := commandCleanRun(cmd)
	if err != nil {
		log.Fatal(err.Error())
	}
}
