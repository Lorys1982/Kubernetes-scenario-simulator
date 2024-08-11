package utils

import (
	"log"
	"main/configs"
	"os"
	"os/exec"
	"strconv"
)

func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
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

func NodeCreate(nodes []configs.Node) {
	for _, node := range nodes {
		nodeName := node.GetName()
		nodeConfName := node.GetConfName()
		replicas := node.GetReplicas()

		input, err := os.ReadFile(nodeConfName)
		if err != nil {
			log.Fatalln(err.Error())
		}

		for i := range replicas {
			fileReplace(nodeConfName, nodeName, nodeName+"-"+strconv.Itoa(i), input...)
			KubectlApply(nodeConfName)
		}
		fileReplace(nodeConfName, nodeName+"-"+strconv.Itoa(replicas), nodeName, input...)
	}
}
