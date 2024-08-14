package utils

import (
	"io"
	"log"
	"main/configs"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// TODO write comment to explain this functions
func commandRun(cmd *exec.Cmd) error {
	outLog, err := os.OpenFile("Logs/stdOut.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	errLog, err := os.OpenFile("Logs/stdErr.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	commandString := strings.Join(cmd.Args, " ")

	_, err = outLog.WriteString("Command: " + commandString + "\n")
	if err != nil {
		return err
	}
	_, err = errLog.WriteString("Command: " + commandString + "\n")
	if err != nil {
		return err
	}

	cmd.Stdout = io.MultiWriter(os.Stdout, outLog)
	cmd.Stderr = io.MultiWriter(os.Stderr, errLog)

	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}

	_, err = outLog.WriteString("\n")
	if err != nil {
		return err
	}
	_, err = errLog.WriteString("\n")
	if err != nil {
		return err
	}
	defer outLog.Close()
	defer errLog.Close()

	return nil
}

// TODO write comment to explain this functions
func commandCleanRun(cmd *exec.Cmd) error {
	outLog, err := os.OpenFile("Logs/stdOut.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	errLog, err := os.OpenFile("Logs/stdErr.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	commandString := strings.Join(cmd.Args, " ")

	_, err = outLog.WriteString("Command: " + commandString + "\n")
	if err != nil {
		return err
	}
	_, err = errLog.WriteString("Command: " + commandString + "\n")
	if err != nil {
		return err
	}

	cmd.Stdout = outLog
	cmd.Stderr = errLog

	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}

	_, err = outLog.WriteString("\n")
	if err != nil {
		return err
	}
	_, err = errLog.WriteString("\n")
	if err != nil {
		return err
	}
	defer outLog.Close()
	defer errLog.Close()

	return nil
}

// clusterArgs function generates args for kwokctl create and kwokctl delete based on given configuration files
//
// # Selectors
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

	if selector && configs.GetSchedulerConf() != "" {
		command = append(command, "--config", configs.GetSchedulerConf())
	}

	if selector && configs.GetAuditConf() != "" {
		command = append(command, "--kube-audit-policy", configs.GetAuditConf())
	}

	return command
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

func KubectlApply(toApply string) {
	args := []string{"apply", "-f", toApply}

	cmd := exec.Command("kubectl", args...)
	err := commandRun(cmd)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// TODO Write comment that ecplains what this does
func NodeCreate(nodes []configs.Node) {
	for _, node := range nodes {
		nodeName := node.GetName()
		nodeConfName := node.GetConfName()
		replicas := node.GetReplicas()
		currentIndex := node.GetCurrentIndex()

		input, err := os.ReadFile(nodeConfName)
		if err != nil {
			log.Fatalln(err.Error())
		}

		for i := range replicas {
			fileReplace(nodeConfName, nodeName, nodeName+"-"+strconv.Itoa(i+currentIndex), input...)
			KubectlApply(nodeConfName)
		}
		// Just restores input (the initial file)
		fileReplace(nodeConfName, "", "", input...)
		node.SetCurrentIndex(replicas)
	}
}
