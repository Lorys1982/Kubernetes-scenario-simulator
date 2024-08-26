package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"main/configs"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

var nodeMutex sync.Mutex
var logMutex sync.Mutex

type Node []configs.Node
type Operations interface {
	Create(float32, int)
	Delete(float32, int)
}

func crashLog(err string) {
	KwokctlDelete()
	errFile, _ := os.OpenFile(fmt.Sprintf("logs/%s_StdErr.log", configs.GetCommandsName()), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	errLog := log.New(errFile, "[Fatal Error] ", log.Ltime|log.Lmicroseconds)
	errLog.Println(err)
	log.Fatal(err)
}

func errLog(err string, cmd string) {
	errFile, _ := os.OpenFile(fmt.Sprintf("logs/%s_StdErr.log", configs.GetCommandsName()), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	errLog := log.New(errFile, "[Error] ", log.Ltime|log.Lmicroseconds)
	errLog.Printf("(Command: %s) %s\n\n", cmd, err)
}

func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func concurrentCommandRun(cmd *exec.Cmd, cfg configs.CommandsList, wg *sync.WaitGroup, cmdSeq int) {
	defer wg.Done()
	_ = commandRun(cmd, cfg.Delay, cmdSeq)
	time.Sleep(time.Duration(cfg.Delay) * time.Second)

}

func concurrentCommandCleanRun(cmd *exec.Cmd, cfg configs.CommandsList, wg *sync.WaitGroup, cmdSeq int) {
	defer wg.Done()
	_ = commandCleanRun(cmd, cfg.Delay, cmdSeq)
	time.Sleep(time.Duration(cfg.Delay) * time.Second)
}

func commandRun(cmd *exec.Cmd, delay float32, cmdSeq ...int) error {
	if len(cmdSeq) == 0 {
		cmdSeq = []int{0}
	}
	outFile, _ := os.OpenFile(fmt.Sprintf("logs/%s_StdOut.log", configs.GetCommandsName()), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	errFile, _ := os.OpenFile(fmt.Sprintf("logs/%s_StdErr.log", configs.GetCommandsName()), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	stdBuff := bufio.NewWriter(outFile)
	errBuff := bufio.NewWriter(errFile)
	defer outFile.Close()
	defer errFile.Close()
	outLog := log.New(stdBuff, fmt.Sprintf("[Command #%d Start] ", cmdSeq[0]), log.Ltime|log.Lmicroseconds)
	errLog := log.New(errBuff, fmt.Sprintf("[Command #%d Start] ", cmdSeq[0]), log.Ltime|log.Lmicroseconds)
	commandString := strings.Join(cmd.Args, " ")

	outLog.Println(commandString)
	errLog.Println(commandString)

	cmd.Stdout = io.MultiWriter(os.Stdout, stdBuff)
	cmd.Stderr = io.MultiWriter(os.Stderr, errBuff)

	_ = cmd.Start()
	cmdErr := cmd.Wait()

	outLog.SetPrefix(fmt.Sprintf("[Command #%d End] ", cmdSeq[0]))
	errLog.SetPrefix(fmt.Sprintf("[Command #%d End] ", cmdSeq[0]))
	outLog.Printf("Delay: %f Seconds\n\n", delay)
	errLog.Printf("Delay: %f Seconds\n\n", delay)

	logMutex.Lock()
	stdBuff.Flush()
	errBuff.Flush()
	logMutex.Unlock()

	return cmdErr
}

func commandCleanRun(cmd *exec.Cmd, delay float32, cmdSeq ...int) error {
	if len(cmdSeq) == 0 {
		cmdSeq = []int{0}
	}
	outFile, _ := os.OpenFile(fmt.Sprintf("logs/%s_StdOut.log", configs.GetCommandsName()), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	errFile, _ := os.OpenFile(fmt.Sprintf("logs/%s_StdErr.log", configs.GetCommandsName()), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	stdBuff := bufio.NewWriter(outFile)
	errBuff := bufio.NewWriter(errFile)
	defer outFile.Close()
	defer errFile.Close()
	outLog := log.New(stdBuff, fmt.Sprintf("[Command #%d Start] ", cmdSeq[0]), log.Ltime|log.Lmicroseconds)
	errLog := log.New(errBuff, fmt.Sprintf("[Command #%d Start] ", cmdSeq[0]), log.Ltime|log.Lmicroseconds)
	commandString := strings.Join(cmd.Args, " ")

	outLog.Println(commandString)
	errLog.Println(commandString)

	cmd.Stdout = stdBuff
	cmd.Stderr = errBuff

	_ = cmd.Start()
	cmdErr := cmd.Wait()

	outLog.SetPrefix(fmt.Sprintf("[Command #%d End] ", cmdSeq[0]))
	errLog.SetPrefix(fmt.Sprintf("[Command #%d End] ", cmdSeq[0]))
	outLog.Printf("Delay: %f Seconds\n\n", delay)
	errLog.Printf("Delay: %f Seconds\n\n", delay)

	logMutex.Lock()
	stdBuff.Flush()
	errBuff.Flush()
	logMutex.Unlock()

	return cmdErr
}

func concurrentExecWrapper(fullCmd []string, cfg configs.CommandsList, wg *sync.WaitGroup, cmdSeq int) {
	defer wg.Done()
	execWrapper(fullCmd, cfg, cmdSeq)
	time.Sleep(time.Duration(cfg.Delay) * time.Second)
}

func execWrapper(fullCmd []string, cfg configs.CommandsList, cmdSeq int) {
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
		object.Create(cfg.Delay, cmdSeq)
		break
	case 0:
		object.Delete(cfg.Delay, cmdSeq)
		break
	default:
		crashLog("No action was provided")
	}
}

func SequentialCommandRun(cmds []configs.CommandsList) {
	var wg sync.WaitGroup
	for cmdSeq, cfg := range cmds {
		cfg.Command = strings.ToLower(cfg.Command)
		if cfg.Concurrent { // The command to run is to be run on a thread
			wg.Add(1)
			if cfg.Exec != "" { // User sent a complete command
				fullCmd := strings.Split(cfg.Exec, " ")
				if !CommandExists(fullCmd[0]) {
					crashLog(fmt.Sprintf("Command %s does not exist", fullCmd[0]))
				}
				cmd := exec.Command(fullCmd[0], fullCmd[1:]...)
				cmd.Dir = "configs/command_configs"
				go concurrentCommandRun(cmd, cfg, &wg, cmdSeq+1)
			} else if cfg.Command != "" { // user sent a wrapped command
				fullCmd := strings.Split(cfg.Command, " ")
				go concurrentExecWrapper(fullCmd, cfg, &wg, cmdSeq+1)
			} else {
				wg.Done()
				crashLog("Invalid Command/Exec")
			}
		} else { // the command to run is to be run sequentially
			wg.Wait()
			if cfg.Exec != "" { // User sent a complete command
				fullCmd := strings.Split(cfg.Exec, " ")
				if !CommandExists(fullCmd[0]) {
					crashLog(fmt.Sprintf("Command %s does not exist", fullCmd[0]))
				}
				cmd := exec.Command(fullCmd[0], fullCmd[1:]...)
				cmd.Dir = "configs/command_configs"
				_ = commandRun(cmd, cfg.Delay, cmdSeq+1)
				fmt.Println("delay of ", cfg.Delay, " seconds")
				time.Sleep(time.Duration(cfg.Delay) * time.Second)
			} else if cfg.Command != "" { // user sent a wrapped command
				fullCmd := strings.Split(cfg.Command, " ")
				execWrapper(fullCmd, cfg, cmdSeq+1)
				fmt.Println("delay of ", cfg.Delay, " seconds")
				time.Sleep(time.Duration(cfg.Delay) * time.Second)
			} else {
				crashLog("Invalid Command/Exec")
			}
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

	if selector && len(configs.GetKwokConf()) != 0 {
		for _, kconf := range configs.GetKwokConf() {
			command = append(command, "--config", kconf)
		}
	}

	if selector && configs.GetAuditConf() != "" {
		command = append(command, "--kube-audit-policy", configs.GetAuditConf())
	}

	return command
}

func KwokctlCreate() {
	args := clusterArgs(true)

	cmd := exec.Command("kwokctl", args...)
	err := commandRun(cmd, 0)
	if err != nil {
		KwokctlDelete()
		crashLog(err.Error())
	}
}

func KwokctlDelete() {
	args := clusterArgs(false)

	cmd := exec.Command("kwokctl", args...)
	err := commandCleanRun(cmd, 0)
	if err != nil {
		return
	}
}

func KubectlApply(toApply string, delay float32, cmdSeq int) {
	args := []string{"apply", "-f", toApply}

	cmd := exec.Command("kubectl", args...)
	_ = commandRun(cmd, delay, cmdSeq)
}

func KubectlDelete(resource string, toDelete string, delay float32, cmdSeq int) {
	args := []string{"delete", resource, toDelete}

	cmd := exec.Command("kubectl", args...)
	_ = commandRun(cmd, delay, cmdSeq)

}

func NodeCreate(nodes Node) {
	nodes.Create(0, 0)
}

func NodeDelete(nodes Node) {
	nodes.Delete(0, 0)
}

func (nodes Node) Create(delay float32, cmdSeq int) {
	for _, node := range nodes {
		nodeName := node.GetName()
		nodeConfName := node.GetConfName()
		replicas := node.GetCount()

		input, err := os.ReadFile(nodeConfName)
		if err != nil {
			crashLog(err.Error())
		}

		for range replicas {
			nodeMutex.Lock()
			currentIndex := node.GetCurrentIndex()
			fileReplace(nodeConfName, nodeName, nodeName+"-"+strconv.Itoa(currentIndex), input...)
			KubectlApply(nodeConfName, delay, cmdSeq)
			node.SetCurrentIndex(currentIndex + 1)
			nodeMutex.Unlock()
		}
		// Just restores input (the initial file)
		fileReplace(nodeConfName, "", "", input...)
	}
}

func (nodes Node) Delete(delay float32, cmdSeq int) {
	for _, node := range nodes {
		nodeName := node.GetName()
		toDelete := node.GetCount()
		initialIndex := node.GetCurrentIndex()

		for range toDelete {
			nodeMutex.Lock()
			currentIndex := node.GetCurrentIndex() - 1
			if currentIndex == -1 {
				nodeMutex.Unlock()
				errLog(fmt.Sprintf("Nodes file: \"%s\"| Nodes name: \"%s\" | Set to delete (count): %d | Deletable: %d",
					node.GetConfName(), nodeName, toDelete, initialIndex),
					"Delete Node")
				break
			}
			KubectlDelete("no", nodeName+"-"+strconv.Itoa(currentIndex), delay, cmdSeq)
			currentIndex--
			node.SetCurrentIndex(currentIndex + 1)
			nodeMutex.Unlock()
		}
	}
}
