package configs

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
)

var conf *Config
var nodeCurrentReplicasVec []nodeCurrentReplicas
var commands *Commands

type Config struct {
	ClusterName string   `yaml:"clusterName"`
	KwokConfigs []string `yaml:"kwokConfigs"`
	Nodes       []Node   `yaml:"nodes"`
	Audit       string   `yaml:"auditLoggingConfig"`
	Commands    string   `yaml:"commandsConfig"`
}

type CommandsList struct {
	Exec       string  `yaml:"exec"`
	Delay      float32 `yaml:"delay"`
	Command    string  `yaml:"command"`
	Filename   string  `yaml:"filename"`
	Count      int     `yaml:"count"`
	Concurrent bool    `yaml:"concurrent"`
}

type Commands struct {
	Kind       string `yaml:"kind"`
	ApiVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec []CommandsList `yaml:"spec"`
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
	ConfigName string `yaml:"name"`
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

func GetCommandsName() string {
	return commands.Metadata.Name
}

func GetCommandsList() []CommandsList {
	return commands.Spec
}

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
	err = yaml.Unmarshal(yamlFile, &commands)
	if err != nil {
		log.Fatal(err.Error())
	}
}
