package utils

import (
	"fmt"
	"io"
	"log"
	"main/configs"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

type Node []configs.Node
type Operations interface {
	Create()
	Delete()
}

func crashLog(err string) {
	KwokctlDelete()
	errLog, _ := os.OpenFile(fmt.Sprintf("logs/%s_StdErr.log", configs.GetCommandsName()), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	errLog.WriteString("Unexpected error: " + err + "\n")
	log.Fatal(err)
}

func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// TODO create the logs directory
func commandRun(cmd *exec.Cmd) error {
	outLog, err := os.OpenFile(fmt.Sprintf("logs/%s_StdOut.log", configs.GetCommandsName()), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	errLog, err := os.OpenFile(fmt.Sprintf("logs/%s_StdErr.log", configs.GetCommandsName()), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
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

// TODO create the logs directory
func commandCleanRun(cmd *exec.Cmd) error {
	outLog, err := os.OpenFile(fmt.Sprintf("logs/%s_StdOut.log", configs.GetCommandsName()), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	errLog, err := os.OpenFile(fmt.Sprintf("logs/%s_StdErr.log", configs.GetCommandsName()), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
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

func execWrapper(fullCmd []string, cfg configs.CommandsList) {
	action := -1
	var object Operations

	for _, cmd := range fullCmd {
		switch cmd {
		case "create":
			action = 1
			break
		case "delete":
			action = 0
			break
		case "node":
			object = Node{
				{
					ConfigName: path.Join("configs/command_configs", cfg.Filename),
					Count:      cfg.Count,
				},
			}
			break
		default:
			crashLog(fmt.Sprintf("Command %s does not exist", cmd))
		}
	}
	if object == nil {
		crashLog("No object provided")
	}

	switch action {
	case 1:
		object.Create()
		break
	case 0:
		object.Delete()
		break
	default:
		crashLog("No action was provided")
	}
}

func SequentialCommandRun(cmds []configs.CommandsList) {
	for _, cfg := range cmds {
		cfg.Command = strings.ToLower(cfg.Command)
		if cfg.Exec != "" { // User sent a complete command
			fullCmd := strings.Split(cfg.Exec, " ")
			if !CommandExists(fullCmd[0]) {
				crashLog(fmt.Sprintf("Command %s does not exist", fullCmd[0]))
			}
			cmd := exec.Command(fullCmd[0], fullCmd[1:]...)
			cmd.Dir = "configs/command_configs"
			err := commandRun(cmd)
			if err != nil {
				crashLog(err.Error())
			}
			fmt.Println("delay of ", cfg.Delay, " seconds")
			time.Sleep(time.Duration(cfg.Delay) * time.Second)
		} else if cfg.Command != "" { // user sent a wrapped command
			fullCmd := strings.Split(cfg.Command, " ")
			execWrapper(fullCmd, cfg)
			fmt.Println("delay of ", cfg.Delay, " seconds")
			time.Sleep(time.Duration(cfg.Delay) * time.Second)
		} else {
			crashLog("Invalid Command/Exec")
		}
	}
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
		crashLog(err.Error())
	}
}

func KwokctlDelete() {
	args := clusterArgs(false)

	cmd := exec.Command("kwokctl", args...)
	err := commandCleanRun(cmd)
	if err != nil {
		return
	}
}

func KubectlApply(toApply string) {
	args := []string{"apply", "-f", toApply}

	cmd := exec.Command("kubectl", args...)
	err := commandRun(cmd)
	if err != nil {
		crashLog(err.Error())
	}
}

func KubectlDelete(resource string, toDelete string) {
	args := []string{"delete", resource, toDelete}

	cmd := exec.Command("kubectl", args...)
	err := commandRun(cmd)
	if err != nil {
		crashLog(err.Error())
	}
}

func NodeCreate(nodes Node) {
	nodes.Create()
}

func NodeDelete(nodes Node) {
	nodes.Delete()
}

func (nodes Node) Create() {
	for _, node := range nodes {
		nodeName := node.GetName()
		nodeConfName := node.GetConfName()
		replicas := node.GetCount()
		currentIndex := node.GetCurrentIndex()

		input, err := os.ReadFile(nodeConfName)
		if err != nil {
			crashLog(err.Error())
		}

		for i := range replicas {
			fileReplace(nodeConfName, nodeName, nodeName+"-"+strconv.Itoa(i+currentIndex), input...)
			KubectlApply(nodeConfName)
		}
		// Just restores input (the initial file)
		fileReplace(nodeConfName, "", "", input...)
		node.SetCurrentIndex(replicas + currentIndex)
	}
}

func (nodes Node) Delete() {
	for _, node := range nodes {
		nodeName := node.GetName()
		toDelete := node.GetCount()
		currentIndex := node.GetCurrentIndex() - 1

		for range toDelete {
			KubectlDelete("no", nodeName+"-"+strconv.Itoa(currentIndex))
			currentIndex--
		}
		node.SetCurrentIndex(currentIndex + 1)
	}
}
