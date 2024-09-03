package configs

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
	"time"
)

var conf *Config
var nodeCurrentReplicasVec []nodeCurrentReplicas
var commandsConf *CommandsConf
var LogTime = time.Now().Format("2006-01-02_15:04:05")
var StartTime time.Time

type Kube struct {
	Filename string
	Args     []string
}

type Config struct {
	ClusterName string   `yaml:"clusterName"`
	KwokConfigs []string `yaml:"kwokConfigs"`
	Nodes       []Node   `yaml:"nodes"`
	Audit       string   `yaml:"auditLoggingConfig"`
	Commands    string   `yaml:"commandsConfig"`
}

type Command struct {
	Exec     string   `yaml:"exec,omitempty"`
	Command  string   `yaml:"command,omitempty"`
	Time     float64  `yaml:"time"`
	Filename string   `yaml:"filename,omitempty"`
	Count    int      `yaml:"count,omitempty"`
	Args     []string `yaml:"args,omitempty"`
	index    int
}

type Queue struct {
	Name       string    `yaml:"name"`
	Kubeconfig string    `yaml:"kubeconfig"`
	Sequence   []Command `yaml:"sequence"`
}

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

type nodeInfo struct {
	Metadata struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
}

type nodeCurrentReplicas struct {
	configName   string
	currentIndex int
}

type Node struct {
	ConfigName string `yaml:"filename"`
	Count      int    `yaml:"count"`
}

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

func (node Node) GetConfName() string {
	return node.ConfigName
}

func (node Node) GetName() string {
	yamlFile, err := os.ReadFile(node.GetConfName())
	var ni *nodeInfo

	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlFile, &ni)
	if err != nil {
		log.Fatal(err)
	}
	return ni.Metadata.Name
}

func (node Node) GetCount() int {
	return node.Count
}

func GetClusterName() string {
	return conf.ClusterName
}

func GetKwokConf() []string {
	return conf.KwokConfigs
}

func GetNodesConf() []Node {
	return conf.Nodes
}

func GetAuditConf() string {
	return conf.Audit
}

func GetCommandsConfName() string { return commandsConf.Metadata.Name }

func GetQueues() []Queue {
	return commandsConf.Spec.Queues
}

func (q Queue) IsEmpty() bool { return q.Name == "" }

func (c Command) GetIndex() int { return c.index }

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

func NewConfig() {
	yamlFile, err := os.ReadFile("configs/config.yaml")

	if err != nil {
		log.Fatal(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatal(err.Error())
	}
	fixFilePath()

	yamlFile, err = os.ReadFile(conf.Commands)

	if err != nil {
		log.Fatal(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &commandsConf)
	if err != nil {
		log.Fatal(err.Error())
	}
	for i := range commandsConf.Spec.Queues {
		for j := range commandsConf.Spec.Queues[i].Sequence {
			commandsConf.Spec.Queues[i].Sequence[j].index = j + 1
		}
	}
}
