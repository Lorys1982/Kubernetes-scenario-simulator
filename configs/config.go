package configs

import (
	"errors"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"time"
)

var conf *Config
var nodeCurrentReplicasVec []nodeCurrentReplicas
var commandsConf *CommandsConf
var LogTime = time.Now().Format("2006-01-02_15:04:05")
var StartTime time.Time

// Option simple implementation from rust, if functions with multiple
// optional parameters are ever needed
type Option[T any] struct {
	none bool
	some T
}

func (o *Option[T]) IsNone() bool {
	return o.none
}
func (o *Option[T]) IsSome() bool {
	return !o.none
}
func (o *Option[T]) None() { o.none = true } // None sets the Option object to none
func (o *Option[T]) Some(data T) {
	o.none = false
	o.some = data
} // Some sets the Option object to some, filling it with data
func (o *Option[T]) GetSome() T {
	return o.some
}

// Kube struct is for wrapped commands who
// want to perform operations through kubectl
type Kube struct {
	Filename string
	Args     []string
}

// Config struct contains all the data of
// the main config file
type Config struct {
	ClusterName string   `yaml:"clusterName"`
	KwokConfigs []string `yaml:"kwokConfigs"`
	Nodes       []Node   `yaml:"nodes"`
	Audit       string   `yaml:"auditLoggingConfig"`
	Commands    string   `yaml:"commandsConfig"`
}

// Command struct contains a single command of the
// command sequence inside a Queue
type Command struct {
	Exec     string   `yaml:"exec,omitempty"`
	Command  string   `yaml:"command,omitempty"`
	Time     float64  `yaml:"time"`
	Filename string   `yaml:"filename,omitempty"`
	Count    int      `yaml:"count,omitempty"`
	Args     []string `yaml:"args,omitempty"`
	index    int
}

// Queue struct contains a Command sequence and the
// data of the Queue
type Queue struct {
	Name       string    `yaml:"name"`
	Kubeconfig string    `yaml:"kubeconfig"`
	Sequence   []Command `yaml:"sequence"`
}

// CommandsConf struct contains the data of the
// commands configuration file
type CommandsConf struct {
	Kind       string `yaml:"kind"`
	ApiVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Aliases []string `yaml:"aliases"`
		Queues  []Queue  `yaml:"queues"`
	}
}

// nodeInfo struct exists only to get the nodes metadata,
// waiting to be expanded
type nodeInfo struct {
	Metadata struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
}

// nodeCurrentReplicas struct is an auxiliary struct to
// keep track of the amount of nodes already deployed
type nodeCurrentReplicas struct {
	configName   string
	currentIndex int
}

// Node struct contains generic nodes information
type Node struct {
	ConfigName string `yaml:"filename"`
	Count      int    `yaml:"count"`
	name       string
}

// GetName returns the name from metadata of the
// corresponding node
func (node Node) GetName() (string, error) {
	yamlFile, err := os.ReadFile(node.GetConfName())
	var ni *nodeInfo

	if err != nil {
		return "", err
	}
	err = yaml.Unmarshal(yamlFile, &ni)
	if err != nil {
		return "", err
	}
	if ni == nil {
		return "", errors.New("Missing parameter: 'name' in " + node.GetConfName())
	}
	return ni.Metadata.Name, nil
}

// GetCurrentIndex fetches the current nodes number from
// the associated node
func (node Node) GetCurrentIndex() int {
	configName := node.ConfigName
	for i := range nodeCurrentReplicasVec {
		if path.Base(nodeCurrentReplicasVec[i].configName) == path.Base(configName) {
			return nodeCurrentReplicasVec[i].currentIndex
		}
	}
	node.SetCurrentIndex(0)
	return 0
}

// SetCurrentIndex sets the number of nodes deployed from the
// associated node
func (node Node) SetCurrentIndex(index int) {
	configName := node.ConfigName
	for i := range nodeCurrentReplicasVec {
		if path.Base(nodeCurrentReplicasVec[i].configName) == path.Base(configName) {
			nodeCurrentReplicasVec[i].currentIndex = index
			return
		}
	}
	nodeCurrentReplicasVec = append(nodeCurrentReplicasVec, nodeCurrentReplicas{
		configName:   configName,
		currentIndex: 0,
	})
}

// GetConfName returns the config file name of the
// associated node
func (node Node) GetConfName() string {
	return node.ConfigName
}

// GetCount returns the replicas wanted for the associated node
func (node Node) GetCount() int {
	return node.Count
}

// GetClusterName returns the cluster's name
func GetClusterName() string {
	return conf.ClusterName
}

// GetKwokConf returns the kwok config file name
func GetKwokConf() []string {
	return conf.KwokConfigs
}

// GetNodesConf returns a list of Node
func GetNodesConf() []Node {
	return conf.Nodes
}

// GetAuditConf returns the audit config file name
func GetAuditConf() string {
	return conf.Audit
}

// GetCommandsConfName returns the name of the commands config file
func GetCommandsConfName() string { return commandsConf.Metadata.Name }

// GetQueues returns a list of Queue
func GetQueues() []Queue {
	return commandsConf.Spec.Queues
}

// IsEmpty returns
//
// true if the associated queue is empty
//
// false otherwise
func (q Queue) IsEmpty() bool { return q.Name == "" }

// GetIndex returns the position in the sequence of the
// associated Command
func (c Command) GetIndex() int { return c.index }

// fixFilePath replaces the filenames in the config file with
// the path to the filename in the topology directory
func fixFilePath() {
	for i, kconf := range conf.KwokConfigs {
		conf.KwokConfigs[i] = path.Join("configs", "topology", kconf)
	}
	if conf.Audit != "" {
		conf.Audit = path.Join("configs", "topology", conf.Audit)
	}
	for i, node := range conf.Nodes {
		conf.Nodes[i].ConfigName = path.Join("configs", "topology", node.ConfigName)
	}
	if conf.Commands != "" {
		conf.Commands = path.Join("configs", "command_configs", conf.Commands)
	}
}

// NewConfig is an initializer, it creates the configuration
// struct from the yaml config files
func NewConfig() {
	yamlFile, err := os.ReadFile("configs/config.yaml")

	if err != nil {
		CrashLog(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		CrashLog(err.Error())
	}
	fixFilePath()

	yamlFile, err = os.ReadFile(conf.Commands)

	if err != nil {
		CrashLog(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &commandsConf)
	if err != nil {
		CrashLog(err.Error())
	}
	for i := range commandsConf.Spec.Queues {
		for j := range commandsConf.Spec.Queues[i].Sequence {
			commandsConf.Spec.Queues[i].Sequence[j].index = j + 1
		}
	}
}
