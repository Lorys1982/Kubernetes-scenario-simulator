package app

import (
	"bytes"
	"fmt"
	"github.com/goaux/decowriter"
	"github.com/notEpsilon/go-pair"
	"io"
	"log"
	"main/configs"
	"main/global"
	. "main/utils"
	"main/writers"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var nodeMutex sync.Mutex
var logChannelStd = writers.LogChannelStd
var logChannelErr = writers.LogChannelErr
var crashLock sync.Mutex
var crashBool = false
var clusterUp = false

type kube configs.Kube
type node []configs.Node
type namespace configs.Namespace

type commandInfo struct {
	QueueName   string
	Kubeconfig  string
	KubeContext configs.Context
	Namespace   string
	CmdIndex    int
	ExecDir     string
}

// Operations interface for resources in command put in configs
//
// float64 is for time to execute, commandInfo is for whatever info on the command is needed
type Operations interface {
	Create(float64, commandInfo)
	Delete(float64, commandInfo)
	Apply(float64, commandInfo)
	Get(float64, commandInfo)
	Scale(float64, commandInfo)
}

// crashHalt function to stop command execution on a fatal error during cluster deletion
// choice:
//   - 1 -> check lock
//   - 0 -> lock
func crashHalt(choice bool) {
	if choice {
		if crashBool {
			crashLock.Lock()
		}
		return
	}
	crashBool = true
	crashLock.Lock()
}

func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func concurrentCommandRun(cmd *exec.Cmd, cmdSpecs configs.Command, wg *sync.WaitGroup, info commandInfo) {
	defer wg.Done()
	time.Sleep(time.Duration(cmdSpecs.Time*float64(time.Second)) * time.Nanosecond)
	log.Printf("Execution at Time: %f\n", cmdSpecs.Time)
	info.ExecDir = "configs/command_configs"
	info.CmdIndex = cmdSpecs.GetIndex()
	crashHalt(true)
	_ = commandRun(cmd, time.Since(global.StartTime).Seconds(), info)
}

func concurrentCommandCleanRun(cmd *exec.Cmd, cmdSpecs configs.Command, wg *sync.WaitGroup, info commandInfo) {
	defer wg.Done()
	time.Sleep(time.Duration(cmdSpecs.Time*float64(time.Second)) * time.Nanosecond)
	log.Printf("Execution at Time: %f\n", cmdSpecs.Time)
	info.ExecDir = "configs/command_configs"
	info.CmdIndex = cmdSpecs.GetIndex()
	crashHalt(true)
	_ = commandCleanRun(cmd, time.Since(global.StartTime).Seconds(), info)
}

func commandRun(cmd *exec.Cmd, execTime float64, info commandInfo) error {
	if info.ExecDir != "" {
		cmd.Dir = info.ExecDir
	}
	stdBuff := bytes.Buffer{}
	errBuff := bytes.Buffer{}
	outLog := log.New(&stdBuff, fmt.Sprintf("[Queue: %s][Command #%d Start] ", info.QueueName, info.CmdIndex), log.Ltime|log.Lmicroseconds)
	errLog := log.New(&errBuff, fmt.Sprintf("[Queue: %s][Command #%d Start] ", info.QueueName, info.CmdIndex), log.Ltime|log.Lmicroseconds)

	commandString := strings.Join(cmd.Args, " ")

	outLog.Println(commandString)
	errLog.Println(commandString)

	stdPrefix := decowriter.New(&stdBuff, []byte(fmt.Sprintf("[Queue: %s][Command #%d] ", info.QueueName, info.CmdIndex)), []byte{})
	errPrefix := decowriter.New(&errBuff, []byte(fmt.Sprintf("[Queue: %s][Command #%d] ", info.QueueName, info.CmdIndex)), []byte{})

	cmd.Stdout = io.MultiWriter(os.Stdout, stdPrefix)
	cmd.Stderr = io.MultiWriter(os.Stderr, errPrefix)

	_ = cmd.Start()
	cmdErr := cmd.Wait()

	r := regexp.MustCompile("(\\r *\\r)*|(\\n *\\n)*")
	stdBuff = *bytes.NewBuffer(r.ReplaceAll(stdBuff.Bytes(), []byte("")))
	errBuff = *bytes.NewBuffer(r.ReplaceAll(errBuff.Bytes(), []byte("")))

	outLog.SetPrefix(fmt.Sprintf("[Queue: %s][Command #%d End] ", info.QueueName, info.CmdIndex))
	errLog.SetPrefix(fmt.Sprintf("[Queue: %s][Command #%d End] ", info.QueueName, info.CmdIndex))
	outLog.Printf("Executed at Time: %f Seconds\n\n", execTime)
	errLog.Printf("Executed at Time: %f Seconds\n\n", execTime)

	logChannelStd <- *pair.New(stdBuff.Bytes(), info.KubeContext.ClusterIndex)
	logChannelErr <- *pair.New(errBuff.Bytes(), info.KubeContext.ClusterIndex)

	return cmdErr
}

func commandCleanRun(cmd *exec.Cmd, execTime float64, info commandInfo) error {
	if info.ExecDir != "" {
		cmd.Dir = info.ExecDir
	}
	stdBuff := bytes.Buffer{}
	errBuff := bytes.Buffer{}
	outLog := log.New(&stdBuff, fmt.Sprintf("[Queue: %s][Command #%d Start] ", info.QueueName, info.CmdIndex), log.Ltime|log.Lmicroseconds)
	errLog := log.New(&errBuff, fmt.Sprintf("[Queue: %s][Command #%d Start] ", info.QueueName, info.CmdIndex), log.Ltime|log.Lmicroseconds)

	commandString := strings.Join(cmd.Args, " ")

	outLog.Println(commandString)
	errLog.Println(commandString)

	stdPrefix := decowriter.New(&stdBuff, []byte(fmt.Sprintf("[Queue: %s][Command #%d] ", info.QueueName, info.CmdIndex)), []byte{})
	errPrefix := decowriter.New(&errBuff, []byte(fmt.Sprintf("[Queue: %s][Command #%d] ", info.QueueName, info.CmdIndex)), []byte{})

	cmd.Stdout = stdPrefix
	cmd.Stderr = errPrefix

	_ = cmd.Start()
	cmdErr := cmd.Wait()

	r := regexp.MustCompile("(\\r *\\r)*|(\\n *\\n)*")
	stdBuff = *bytes.NewBuffer(r.ReplaceAll(stdBuff.Bytes(), []byte("")))
	errBuff = *bytes.NewBuffer(r.ReplaceAll(errBuff.Bytes(), []byte("")))

	outLog.SetPrefix(fmt.Sprintf("[Queue: %s][Command #%d End] ", info.QueueName, info.CmdIndex))
	errLog.SetPrefix(fmt.Sprintf("[Queue: %s][Command #%d End] ", info.QueueName, info.CmdIndex))
	outLog.Printf("Executed at Time: %f Seconds\n\n", execTime)
	errLog.Printf("Executed at Time: %f Seconds\n\n", execTime)

	logChannelStd <- *pair.New(stdBuff.Bytes(), info.KubeContext.ClusterIndex)
	logChannelErr <- *pair.New(errBuff.Bytes(), info.KubeContext.ClusterIndex)

	return cmdErr
}

func concurrentExecWrapper(fullCmd []string, cmdSpecs configs.Command, wg *sync.WaitGroup, info commandInfo) {
	defer wg.Done()
	time.Sleep(time.Duration(cmdSpecs.Time*float64(time.Second)) * time.Nanosecond)
	log.Printf("Execution at Time: %f\n", cmdSpecs.Time)
	info.ExecDir = "configs/command_configs"
	info.CmdIndex = cmdSpecs.GetIndex()
	crashHalt(true)
	execWrapper(fullCmd, cmdSpecs, info)
}

func execWrapper(fullCmd []string, cmdSpecs configs.Command, info commandInfo) {
	var object Operations
	resource := fullCmd[0]
	action := fullCmd[1]
	//params := fullCmd[2:]

	// Resource switch
	switch resource {
	case "node":
		object = node{
			{
				Filename: path.Join("configs/command_configs", cmdSpecs.Filename),
				Count:    cmdSpecs.Count,
			},
		}
	case "kube":
		object = kube{
			Filename: cmdSpecs.Filename,
			Resource: cmdSpecs.Resource,
			Args:     cmdSpecs.Args,
			Count:    cmdSpecs.Count,
		}
	case "namespace":
		object = namespace{
			Args: cmdSpecs.Args,
		}
	default:
		crashLog(fmt.Sprintf("Resource %s does not exist", fullCmd[0]), info)
	}

	// Action switch
	switch action {
	case "create":
		object.Create(time.Since(global.StartTime).Seconds(), info)
		break
	case "delete":
		object.Delete(time.Since(global.StartTime).Seconds(), info)
		break
	case "get":
		object.Get(time.Since(global.StartTime).Seconds(), info)
	case "apply":
		object.Apply(time.Since(global.StartTime).Seconds(), info)
	case "scale":
		object.Scale(time.Since(global.StartTime).Seconds(), info)
	default:
		crashLog("No action was provided", info)
	}
}

func ConcurrentQueueRun(clusterQueues [][]configs.Queue) {
	var wgQueues sync.WaitGroup
	global.StartTime = time.Now()
	for _, queues := range clusterQueues {
		for _, queue := range queues {
			wgQueues.Add(1)
			go ConcurrentCommandsRun(queue, &wgQueues)
		}
	}
	wgQueues.Wait()
}

func ConcurrentCommandsRun(queue configs.Queue, wgQueues *sync.WaitGroup) {
	wgCommands := sync.WaitGroup{}
	cmds := queue.Sequence
	sort.Slice(cmds, func(i, j int) bool {
		return cmds[j].Time > cmds[i].Time
	})

	info := commandInfo{
		QueueName:   queue.Name,
		Kubeconfig:  queue.Kubeconfig,
		KubeContext: queue.KubeContext,
		Namespace:   "",
		CmdIndex:    -1,
		ExecDir:     "",
	}

	for _, cmdSpecs := range cmds {
		wgCommands.Add(1)
		cmdSpecs.Command = strings.ToLower(cmdSpecs.Command)
		if len(cmdSpecs.Context) != 0 {
			info.KubeContext.Name = cmdSpecs.Context
		}
		if len(cmdSpecs.Namespace) != 0 {
			info.Namespace = cmdSpecs.Namespace
		}
		if cmdSpecs.Exec != "" { // User sent a complete command
			fullCmd := strings.Split(cmdSpecs.Exec, " ")
			if !CommandExists(fullCmd[0]) {
				info.CmdIndex = cmdSpecs.GetIndex()
				crashLog(fmt.Sprintf("Command %s does not exist", fullCmd[0]), info)
			}
			fullCmd = kubectlScopeArgs(fullCmd, info)
			cmd := exec.Command(fullCmd[0], fullCmd[1:]...)
			cmd.Dir = "configs/command_configs"
			go concurrentCommandRun(cmd, cmdSpecs, &wgCommands, info)
		} else if cmdSpecs.Command != "" { // user sent a wrapped command
			fullCmd := strings.Split(cmdSpecs.Command, " ")
			go concurrentExecWrapper(fullCmd, cmdSpecs, &wgCommands, info)
		} else {
			info.CmdIndex = cmdSpecs.GetIndex()
			crashLog("Invalid Command/Exec", info)
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
func clusterArgs(selector bool) [][]string {
	var command [][]string
	clusters := configs.GetClusterNames()
	kwokConf := configs.GetKwokConf()
	auditConf := configs.GetAuditConf()

	if selector {
		for range clusters {
			command = append(command, []string{"create", "cluster"})
		}
	} else {
		for range clusters {
			command = append(command, []string{"delete", "cluster"})
		}
	}

	for i := range clusters {
		if clusters[i] != "" {
			command[i] = append(command[i], "--name", clusters[i])
		}

		if selector && len(kwokConf) != 0 {
			for _, kconf := range kwokConf[i] {
				command[i] = append(command[i], kconf)
			}
		}

		if selector && auditConf[i] != "" {
			command[i] = append(command[i], "--kube-audit-policy", auditConf[i])
		}
		if selector && configs.IsLiqoActive() {
			command[i] = append(command[i], "--runtime", "kind")
		}
	}

	return command
}

func KwokctlCreateAll() {
	wg := sync.WaitGroup{}
	args := clusterArgs(true)
	for i := range args {
		wg.Add(1)
		info := commandInfo{
			QueueName: "KwokctlCreate",
			KubeContext: configs.Context{
				ClusterIndex: i,
			},
			CmdIndex: 0,
			ExecDir:  "configs/topology",
		}
		// IMPORTANT! Until a PR is acceptend by KWOK, this concurrency section only works with a modified
		// kwokctl binary (provided in the application), if using the standard kwokctl, please comment the if
		// statement below and only keep what is in the else clause
		// plus, to use more than 2 clusters we need to run this command
		// sudo sysctl fs.inotify.max_user_watches=524288
		// sudo sysctl fs.inotify.max_user_instances=512
		if configs.IsLiqoActive() {
			go kwokctlCreate(args[i], info, &wg)
			nodeName := "kwok" + "-" + configs.GetClusterName(i) + "-" + "control-plane"
			err := WaitForContainer(nodeName)
			if err != nil {
				crashLog(err.Error(), info)
			}
		} else {
			kwokctlCreate(args[i], info, &wg)
		}
	}
	wg.Wait()
	clusterUp = true
}

func kwokctlCreate(args []string, info commandInfo, wg *sync.WaitGroup) {
	cmd := exec.Command("kwokctl", args...)
	err := commandRun(cmd, 0, info)
	if err != nil {
		crashLog(err.Error(), info)
	}
	wg.Done()
}

func KwokctlDeleteAll() {
	wg := sync.WaitGroup{}
	args := clusterArgs(false)
	clusters := configs.GetClusterNames()
	crashHalt(false)
	for i := range args {
		wg.Add(1)
		go kwokctlDelete(args[i], clusters[i], i, &wg)
	}
	wg.Wait()
	writers.KillLoggers()
}

func kwokctlDelete(args []string, cluster string, clusterIndex int, wg *sync.WaitGroup) {
	home, _ := os.UserHomeDir()

	// Copy and compress log file
	Compress("audit.log", path.Join(home, ".kwok/clusters", cluster, "logs"), clusterIndex)

	cmd := exec.Command("kwokctl", args...)
	err := commandCleanRun(cmd, time.Since(global.StartTime).Seconds(), commandInfo{
		QueueName: "KwokctlDelete",
		KubeContext: configs.Context{
			ClusterIndex: clusterIndex,
		},
		CmdIndex: 0,
		ExecDir:  "configs/topology",
	})
	if err != nil {
		log.Fatal("Kwokctl cluster deletion failed - ", err)
	}
	wg.Done()
}

func KubectlApply(toApply string, execTime float64, info commandInfo, cmdArgs ...string) {
	args := []string{"apply", "-f", toApply}
	args = append(args, cmdArgs...)
	args = kubectlScopeArgs(args, info)

	cmd := exec.Command("kubectl", args...)
	_ = commandRun(cmd, execTime, info)
}

func KubectlCreate(resource string, execTime float64, info commandInfo, cmdArgs ...string) {
	args := []string{"create", resource}
	args = append(args, cmdArgs...)
	args = kubectlScopeArgs(args, info)

	cmd := exec.Command("kubectl", args...)
	_ = commandRun(cmd, execTime, info)
}

func KubectlDelete(resource string, execTime float64, info commandInfo, cmdArgs ...string) {
	args := []string{"delete", resource}
	args = append(args, cmdArgs...)
	args = kubectlScopeArgs(args, info)

	cmd := exec.Command("kubectl", args...)
	_ = commandRun(cmd, execTime, info)

}

func KubectlGet(resource string, execTime float64, info commandInfo, cmdArgs ...string) {
	args := []string{"get", resource}
	args = append(args, cmdArgs...)
	args = kubectlScopeArgs(args, info)

	cmd := exec.Command("kubectl", args...)
	_ = commandRun(cmd, execTime, info)
}

func kubectScale(resource string, replicas int, execTime float64, info commandInfo, cmdArgs ...string) {
	args := []string{"scale", resource}
	args = append(args, cmdArgs...)
	args = append(args, "--replicas", strconv.Itoa(replicas))
	args = kubectlScopeArgs(args, info)

	cmd := exec.Command("kubectl", args...)
	_ = commandRun(cmd, execTime, info)
}

func KubectlCordon(execTime float64, info commandInfo, cmdArgs ...string) {
	args := []string{"cordon"}
	args = append(args, cmdArgs...)
	args = kubectlScopeArgs(args, info)

	cmd := exec.Command("kubectl", args...)
	_ = commandRun(cmd, execTime, info)
}

func KubectlUncordon(execTime float64, info commandInfo, cmdArgs ...string) {
	args := []string{"uncordon"}
	args = append(args, cmdArgs...)
	args = kubectlScopeArgs(args, info)

	cmd := exec.Command("kubectl", args...)
	_ = commandRun(cmd, execTime, info)
}

func kubectlScopeArgs(args []string, info commandInfo) []string {
	args = append(args, "--kubeconfig", info.Kubeconfig)
	args = append(args, "--context", info.KubeContext.Name)
	if len(info.Namespace) != 0 {
		if info.Namespace == "all" {
			args = append(args, "--all-namespaces")
		} else {
			args = append(args, "--namespace", info.Namespace)
		}
	}
	return args
}

func NodeCreate(nodes node, clusterIndex int) {
	nodes.Create(0, commandInfo{
		QueueName:   "Topology",
		Kubeconfig:  configs.GetKubeConfigPath(clusterIndex),
		KubeContext: configs.ClusterKubeconfigs[clusterIndex].Contexts[0],
		CmdIndex:    0,
		ExecDir:     "configs/topology",
	})
}

func NodeDelete(nodes node, clusterIndex int) {
	nodes.Delete(0, commandInfo{
		QueueName:   "Topology",
		Kubeconfig:  configs.GetKubeConfigPath(clusterIndex),
		KubeContext: configs.ClusterKubeconfigs[clusterIndex].Contexts[0],
		CmdIndex:    0,
		ExecDir:     "configs/topology",
	})
}

func (nodes node) Create(execTime float64, info commandInfo) {
	for _, node := range nodes {
		nodeName, err := node.GetName()
		if err != nil {
			crashLog(err.Error(), info)
		}
		nodeConfName := node.GetConfName()
		replicas := node.GetCount()

		input, err := os.ReadFile(nodeConfName)
		if err != nil {
			crashLog(err.Error(), info)
		}
		opt := Some(input)

		for range replicas {
			nodeMutex.Lock()
			currentIndex := node.GetCurrentIndex(info.KubeContext.ClusterIndex)
			FileReplace(nodeConfName, nodeName, nodeName+"-"+strconv.Itoa(currentIndex), opt)
			KubectlApply(path.Base(nodeConfName), execTime, info)
			node.SetCurrentIndex(currentIndex+1, info.KubeContext.ClusterIndex)
			nodeMutex.Unlock()
		}
		// Just restores input (the initial file)
		FileReplace(nodeConfName, "", "", opt)
	}
}

func (nodes node) Delete(execTime float64, info commandInfo) {
	for _, node := range nodes {
		nodeName, err := node.GetName()
		if err != nil {
			crashLog(err.Error(), info)
		}
		toDelete := node.GetCount()
		initialIndex := node.GetCurrentIndex(info.KubeContext.ClusterIndex)

		for range toDelete {
			nodeMutex.Lock()
			currentIndex := node.GetCurrentIndex(info.KubeContext.ClusterIndex) - 1
			if currentIndex == -1 {
				nodeMutex.Unlock()
				errLog(fmt.Sprintf("Nodes file: \"%s\" | Nodes name: \"%s\" | Set to delete (count): %d | Deletable: %d",
					node.GetConfName(), nodeName, toDelete, initialIndex),
					"Delete node", info)
				break
			}
			KubectlDelete("no", execTime, info, []string{nodeName + "-" + strconv.Itoa(currentIndex)}...)
			currentIndex--
			node.SetCurrentIndex(currentIndex+1, info.KubeContext.ClusterIndex)
			nodeMutex.Unlock()
		}
	}
}

func (nodes node) Apply(execTime float64, info commandInfo) {
	nodes.Create(execTime, info)
}

func (nodes node) Get(execTime float64, info commandInfo) {
	KubectlGet("no", execTime, info)
}

func (nodes node) Scale(execTime float64, info commandInfo) {
	wantedReplicas := nodes[0].GetCount()
	currentReplicas := nodes[0].GetCurrentIndex(info.KubeContext.ClusterIndex)

	if currentReplicas > wantedReplicas {
		nodes[0].Count = currentReplicas - wantedReplicas
		nodes.Delete(execTime, info)
	} else if currentReplicas < wantedReplicas {
		nodes[0].Count = wantedReplicas - currentReplicas
		nodes.Create(execTime, info)
	}
}

func (k kube) Create(execTime float64, info commandInfo) {
	k.Args = fixArgs(k.Args)
	KubectlCreate(k.Resource, execTime, info, k.Args...)
}

func (k kube) Delete(execTime float64, info commandInfo) {
	k.Args = fixArgs(k.Args)
	if k.Filename != "" {
		k.Args = append([]string{"-f", k.Filename}, k.Args...)
	}
	KubectlDelete(k.Resource, execTime, info, k.Args...)
}

func (k kube) Apply(execTime float64, info commandInfo) {
	k.Args = fixArgs(k.Args)
	KubectlApply(k.Filename, execTime, info, k.Args...)
}

func (k kube) Get(execTime float64, info commandInfo) {
	k.Args = fixArgs(k.Args)
	KubectlGet(k.Resource, execTime, info, k.Args...)
}

func (k kube) Scale(execTime float64, info commandInfo) {
	k.Args = fixArgs(k.Args)
	if k.Filename != "" {
		k.Args = append([]string{"-f", k.Filename}, k.Args...)
	}
	kubectScale(k.Resource, k.Count, execTime, info, k.Args...)
}

func (n namespace) Create(execTime float64, info commandInfo) {
	n.Args = fixArgs(n.Args)
	KubectlCreate("namespace", execTime, info, n.Args...)
}

func (n namespace) Delete(execTime float64, info commandInfo) {
	n.Args = fixArgs(n.Args)
	KubectlDelete("namespace", execTime, info, n.Args...)
}

func (n namespace) Apply(execTime float64, info commandInfo) {
	panic("This Function should never be called")
}

func (n namespace) Get(execTime float64, info commandInfo) {
	n.Args = fixArgs(n.Args)
	KubectlGet("namespace", execTime, info, n.Args...)
}

func (n namespace) Scale(execTime float64, info commandInfo) {
	panic("This Function should never be called")
}

func fixArgs(args []string) []string {
	var rArgs []string
	for _, arg := range args {
		rArgs = append(rArgs, strings.Split(arg, " ")...)
	}
	return rArgs
}

func errLog(err string, s string, info commandInfo) {
	logInfo := global.LogCommandInfo{
		CmdIndex:     info.CmdIndex,
		QueueName:    info.QueueName,
		ClusterIndex: info.KubeContext.ClusterIndex,
	}
	writers.ErrLog(err, s, logInfo)
}

func crashLog(err string, info commandInfo) {
	logInfo := global.LogCommandInfo{
		CmdIndex:     info.CmdIndex,
		QueueName:    info.QueueName,
		ClusterIndex: info.KubeContext.ClusterIndex,
	}
	if clusterUp {
		writers.CrashLog(err, Some(logInfo), KwokctlDeleteAll)
	} else {
		writers.CrashLog(err, Some(logInfo), KwokctlDeleteAll)
	}
}
