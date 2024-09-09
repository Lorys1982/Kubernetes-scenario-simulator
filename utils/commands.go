package utils

import (
	"bufio"
	"fmt"
	"github.com/goaux/decowriter"
	"io"
	"log"
	"main/configs"
	"os"
	"os/exec"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var nodeMutex sync.Mutex
var logMutex sync.Mutex

type kube configs.Kube
type node []configs.Node

// Operations interface for resources in command put in configs
//
// float64 is for time to execute, commandInfo is for whatever info on the command is needed
type Operations interface {
	Create(float64, commandInfo)
	Delete(float64, commandInfo)
	Apply(float64, commandInfo)
	Get(float64, commandInfo)
}

type commandInfo struct {
	Queue  configs.Queue
	CmdSeq int
}

func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func concurrentCommandRun(cmd *exec.Cmd, cfg configs.Command, wg *sync.WaitGroup, queue configs.Queue) {
	defer wg.Done()
	time.Sleep(time.Duration(cfg.Time) * time.Second)
	log.Printf("Execution at Time: %f\n", cfg.Time)
	info := commandInfo{
		Queue:  queue,
		CmdSeq: cfg.GetIndex(),
	}
	_ = commandRun(cmd, time.Since(configs.StartTime).Seconds(), info)
}

func concurrentCommandCleanRun(cmd *exec.Cmd, cfg configs.Command, wg *sync.WaitGroup, queue configs.Queue) {
	defer wg.Done()
	time.Sleep(time.Duration(cfg.Time) * time.Second)
	log.Printf("Execution at Time: %f\n", cfg.Time)
	info := commandInfo{
		Queue:  queue,
		CmdSeq: cfg.GetIndex(),
	}
	_ = commandCleanRun(cmd, time.Since(configs.StartTime).Seconds(), info)
}

func commandRun(cmd *exec.Cmd, execTime float64, info commandInfo) error {
	if info.Queue.IsEmpty() {
		info.Queue.Name = "<None>"
	}
	outFile, _ := os.OpenFile(fmt.Sprintf("logs/%s_StdOut_%s.log", configs.GetCommandsConfName(), configs.LogTime), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	errFile, _ := os.OpenFile(fmt.Sprintf("logs/%s_StdErr_%s.log", configs.GetCommandsConfName(), configs.LogTime), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	stdBuff := bufio.NewWriter(outFile)
	errBuff := bufio.NewWriter(errFile)
	defer outFile.Close()
	defer errFile.Close()
	outLog := log.New(stdBuff, fmt.Sprintf("[Queue: %s][Command #%d Start] ", info.Queue.Name, info.CmdSeq), log.Ltime|log.Lmicroseconds)
	errLog := log.New(errBuff, fmt.Sprintf("[Queue: %s][Command #%d Start] ", info.Queue.Name, info.CmdSeq), log.Ltime|log.Lmicroseconds)

	commandString := strings.Join(cmd.Args, " ")

	outLog.Println(commandString)
	errLog.Println(commandString)

	stdPrefix := decowriter.New(stdBuff, []byte(fmt.Sprintf("[Queue: %s][Command #%d] ", info.Queue.Name, info.CmdSeq)), []byte{})
	errPrefix := decowriter.New(errBuff, []byte(fmt.Sprintf("[Queue: %s][Command #%d] ", info.Queue.Name, info.CmdSeq)), []byte{})

	cmd.Stdout = io.MultiWriter(os.Stdout, stdPrefix)
	cmd.Stderr = io.MultiWriter(os.Stderr, errPrefix)

	_ = cmd.Start()
	cmdErr := cmd.Wait()

	outLog.SetPrefix(fmt.Sprintf("[Queue: %s][Command #%d End] ", info.Queue.Name, info.CmdSeq))
	errLog.SetPrefix(fmt.Sprintf("[Queue: %s][Command #%d End] ", info.Queue.Name, info.CmdSeq))
	outLog.Printf("Executed at Time: %f Seconds\n\n", execTime)
	errLog.Printf("Executed at Time: %f Seconds\n\n", execTime)

	logMutex.Lock()
	stdBuff.Flush()
	errBuff.Flush()
	logMutex.Unlock()

	return cmdErr
}

func commandCleanRun(cmd *exec.Cmd, execTime float64, info commandInfo) error {
	if info.Queue.IsEmpty() {
		info.Queue.Name = "<None>"
	}
	outFile, _ := os.OpenFile(fmt.Sprintf("logs/%s_StdOut_%s.log", configs.GetCommandsConfName(), configs.LogTime), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	errFile, _ := os.OpenFile(fmt.Sprintf("logs/%s_StdErr_%s.log", configs.GetCommandsConfName(), configs.LogTime), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	stdBuff := bufio.NewWriter(outFile)
	errBuff := bufio.NewWriter(errFile)
	defer outFile.Close()
	defer errFile.Close()
	outLog := log.New(stdBuff, fmt.Sprintf("[Queue: %s][Command #%d Start] ", info.Queue.Name, info.CmdSeq), log.Ltime|log.Lmicroseconds)
	errLog := log.New(errBuff, fmt.Sprintf("[Queue: %s][Command #%d Start] ", info.Queue.Name, info.CmdSeq), log.Ltime|log.Lmicroseconds)

	commandString := strings.Join(cmd.Args, " ")

	outLog.Println(commandString)
	errLog.Println(commandString)

	stdPrefix := decowriter.New(stdBuff, []byte(fmt.Sprintf("[Queue: %s][Command #%d] ", info.Queue.Name, info.CmdSeq)), []byte{})
	errPrefix := decowriter.New(errBuff, []byte(fmt.Sprintf("[Queue: %s][Command #%d] ", info.Queue.Name, info.CmdSeq)), []byte{})

	cmd.Stdout = stdPrefix
	cmd.Stderr = errPrefix

	_ = cmd.Start()
	cmdErr := cmd.Wait()

	outLog.SetPrefix(fmt.Sprintf("[Queue: %s][Command #%d End] ", info.Queue.Name, info.CmdSeq))
	errLog.SetPrefix(fmt.Sprintf("[Queue: %s][Command #%d End] ", info.Queue.Name, info.CmdSeq))
	outLog.Printf("Executed at Time: %f Seconds\n\n", execTime)
	errLog.Printf("Executed at Time: %f Seconds\n\n", execTime)

	logMutex.Lock()
	stdBuff.Flush()
	errBuff.Flush()
	logMutex.Unlock()

	return cmdErr
}

func concurrentExecWrapper(fullCmd []string, cfg configs.Command, wg *sync.WaitGroup, queue configs.Queue) {
	defer wg.Done()
	time.Sleep(time.Duration(cfg.Time) * time.Second)
	log.Printf("Execution at Time: %f\n", cfg.Time)
	info := commandInfo{
		Queue:  queue,
		CmdSeq: cfg.GetIndex(),
	}
	execWrapper(fullCmd, cfg, info)
}

func execWrapper(fullCmd []string, cfg configs.Command, info commandInfo) {
	var object Operations
	resource := fullCmd[0]
	action := fullCmd[1]
	//params := fullCmd[2:]

	// Resource switch
	switch resource {
	case "node":
		object = node{
			{
				ConfigName: path.Join("configs/command_configs", cfg.Filename),
				Count:      cfg.Count,
			},
		}
	case "kube":
		object = kube{
			Filename: path.Join("configs/command_configs", cfg.Filename),
			Args:     cfg.Args,
		}
	default:
		crashLog(fmt.Sprintf("Resource %s does not exist", fullCmd[0]))
	}

	// Action switch
	switch action {
	case "create":
		object.Create(time.Since(configs.StartTime).Seconds(), info)
		break
	case "delete":
		object.Delete(time.Since(configs.StartTime).Seconds(), info)
		break
	case "get":
		object.Get(time.Since(configs.StartTime).Seconds(), info)
	case "apply":
		object.Apply(time.Since(configs.StartTime).Seconds(), info)
	default:
		crashLog("No action was provided")
	}
}

func ConcurrentQueueRun(queues []configs.Queue) {
	var wgQueues sync.WaitGroup
	configs.StartTime = time.Now()
	for _, queue := range queues {
		wgQueues.Add(1)
		go ConcurrentCommandsRun(queue, &wgQueues)
	}
	wgQueues.Wait()
}

func ConcurrentCommandsRun(queue configs.Queue, wgQueues *sync.WaitGroup) {
	wgCommands := sync.WaitGroup{}
	cmds := queue.Sequence
	sort.Slice(cmds, func(i, j int) bool {
		return cmds[j].Time > cmds[i].Time
	})
	for _, cfg := range cmds {
		wgCommands.Add(1)
		cfg.Command = strings.ToLower(cfg.Command)
		if cfg.Exec != "" { // User sent a complete command
			fullCmd := strings.Split(cfg.Exec, " ")
			if !CommandExists(fullCmd[0]) {
				crashLog(fmt.Sprintf("Command %s does not exist", fullCmd[0]))
			}
			if queue.Kubeconfig != "" {
				fullCmd = append(fullCmd, "--kubeconfig", queue.Kubeconfig)
			}
			cmd := exec.Command(fullCmd[0], fullCmd[1:]...)
			cmd.Dir = "configs/command_configs"
			go concurrentCommandRun(cmd, cfg, &wgCommands, queue)
		} else if cfg.Command != "" { // user sent a wrapped command
			fullCmd := strings.Split(cfg.Command, " ")
			if queue.Kubeconfig != "" {
				fullCmd = append(fullCmd, "--kubeconfig", queue.Kubeconfig)
			}
			go concurrentExecWrapper(fullCmd, cfg, &wgCommands, queue)
		} else {
			crashLog("Invalid Command/Exec")
		}
	}
	wgCommands.Wait()
	wgQueues.Done()
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
	err := commandRun(cmd, 0, commandInfo{
		Queue:  configs.Queue{},
		CmdSeq: 0,
	})
	if err != nil {
		KwokctlDelete()
		crashLog(err.Error())
	}
}

func KwokctlDelete() {
	args := clusterArgs(false)

	cmd := exec.Command("kwokctl", args...)
	err := commandCleanRun(cmd, time.Since(configs.StartTime).Seconds(), commandInfo{
		Queue:  configs.Queue{},
		CmdSeq: 0,
	})
	if err != nil {
		return
	}
}

func KubectlApply(toApply string, execTime float64, info commandInfo, cmdArgs ...string) {
	args := []string{"apply", "-f", toApply}
	args = append(args, cmdArgs...)
	if info.Queue.Kubeconfig != "" {
		args = append(args, "--kubeconfig", info.Queue.Kubeconfig)
	}

	cmd := exec.Command("kubectl", args...)
	_ = commandRun(cmd, execTime, info)
}

func KubectlCreate(toCreate string, execTime float64, info commandInfo, cmdArgs ...string) {
	args := []string{"create", "-f", toCreate}
	args = append(args, cmdArgs...)
	if info.Queue.Kubeconfig != "" {
		args = append(args, "--kubeconfig", info.Queue.Kubeconfig)
	}

	cmd := exec.Command("kubectl", args...)
	_ = commandRun(cmd, execTime, info)
}

func KubectlDelete(execTime float64, info commandInfo, cmdArgs ...string) {
	args := []string{"delete"}
	args = append(args, cmdArgs...)
	if info.Queue.Kubeconfig != "" {
		args = append(args, "--kubeconfig", info.Queue.Kubeconfig)
	}

	cmd := exec.Command("kubectl", args...)
	_ = commandRun(cmd, execTime, info)

}

func KubectlGet(execTime float64, info commandInfo, cmdArgs ...string) {
	args := []string{"get"}
	args = append(args, cmdArgs...)
	if info.Queue.Kubeconfig != "" {
		args = append(args, "--kubeconfig", info.Queue.Kubeconfig)
	}

	cmd := exec.Command("kubectl", args...)
	_ = commandRun(cmd, execTime, info)
}

func NodeCreate(nodes node) {
	nodes.Create(0, commandInfo{
		Queue:  configs.Queue{},
		CmdSeq: 0,
	})
}

func NodeDelete(nodes node) {
	nodes.Delete(0, commandInfo{
		Queue:  configs.Queue{},
		CmdSeq: 0,
	})
}

func (nodes node) Create(execTime float64, info commandInfo) {
	var optionInput = &configs.Option[[]byte]{}
	for _, node := range nodes {
		nodeName, err := node.GetName()
		if err != nil {
			crashLog(err.Error())
		}
		nodeConfName := node.GetConfName()
		replicas := node.GetCount()

		input, err := os.ReadFile(nodeConfName)
		if err != nil {
			crashLog(err.Error())
		}
		optionInput.Some(input)

		for range replicas {
			nodeMutex.Lock()
			currentIndex := node.GetCurrentIndex()
			fileReplace(nodeConfName, nodeName, nodeName+"-"+strconv.Itoa(currentIndex), *optionInput)
			KubectlApply(nodeConfName, execTime, info)
			node.SetCurrentIndex(currentIndex + 1)
			nodeMutex.Unlock()
		}
		// Just restores input (the initial file)
		fileReplace(nodeConfName, "", "", *optionInput)
	}
}

func (nodes node) Delete(execTime float64, info commandInfo) {
	for _, node := range nodes {
		nodeName, err := node.GetName()
		if err != nil {
			crashLog(err.Error())
		}
		toDelete := node.GetCount()
		initialIndex := node.GetCurrentIndex()

		for range toDelete {
			nodeMutex.Lock()
			currentIndex := node.GetCurrentIndex() - 1
			if currentIndex == -1 {
				nodeMutex.Unlock()
				errLog(fmt.Sprintf("Nodes file: \"%s\" | Nodes name: \"%s\" | Set to delete (count): %d | Deletable: %d",
					node.GetConfName(), nodeName, toDelete, initialIndex),
					"Delete node")
				break
			}
			KubectlDelete(execTime, info, []string{"no", nodeName + "-" + strconv.Itoa(currentIndex)}...)
			currentIndex--
			node.SetCurrentIndex(currentIndex + 1)
			nodeMutex.Unlock()
		}
	}
}

func (nodes node) Apply(execTime float64, info commandInfo) {
	nodes.Create(execTime, info)
}

func (nodes node) Get(execTime float64, info commandInfo) {
	KubectlGet(execTime, info, "no")
}

func (k kube) Create(execTime float64, info commandInfo) {
	k.Args = fixArgs(k.Args)
	KubectlCreate(k.Filename, execTime, info, k.Args...)
}

func (k kube) Delete(execTime float64, info commandInfo) {
	k.Args = fixArgs(k.Args)
	if k.Filename != "" {
		k.Args = append([]string{"-f", k.Filename}, k.Args...)
	}
	KubectlDelete(execTime, info, k.Args...)
}

func (k kube) Apply(execTime float64, info commandInfo) {
	k.Args = fixArgs(k.Args)
	KubectlApply(k.Filename, execTime, info, k.Args...)
}

func (k kube) Get(execTime float64, info commandInfo) {
	k.Args = fixArgs(k.Args)
	KubectlGet(execTime, info, k.Args...)
}

func fixArgs(args []string) []string {
	var rArgs []string
	for _, arg := range args {
		rArgs = append(rArgs, strings.Split(arg, " ")...)
	}
	return rArgs
}

func errLog(err string, s string) {
	configs.ErrLog(err, s)
}

func crashLog(err string) {
	KwokctlDelete()
	configs.CrashLog(err)
}
