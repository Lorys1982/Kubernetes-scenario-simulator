package v1alpha1

import (
	"errors"
	"gopkg.in/yaml.v2"
	"os"
)

const (
	ResourceKindSimConfiguration = "SimConfiguration"
)

var conf *Config
var nodeCurrentReplicasVec [][]nodeCurrentReplicas

// Config struct contains all the data in
// the main config file about the topology of the cluster(s)
type Config struct {
	Kind              string    `yaml:"kind"`
	ApiVersion        string    `yaml:"apiVersion"`
	Liqo              LiqoOpt   `yaml:"liqo"`
	Clusters          []Cluster `yaml:"clusters"`
	LiqoConsumerIndex int       `yaml:"-"`
}

type Cluster struct {
	ClusterName string   `yaml:"clusterName"`
	KwokConfigs []string `yaml:"kwokConfigs"`
	Nodes       []Node   `yaml:"nodes"`
	Audit       string   `yaml:"auditLoggingConfig"`
	Commands    string   `yaml:"commandsConfig"`
}

type LiqoOffload struct {
	Namespace         string   `yaml:"namespace"`
	ClusterSelector   []string `yaml:"clusterSelector"`   // Default: empty
	NamespaceStrategy string   `yaml:"namespaceStrategy"` // Default: DefaultName
	PodStrategy       string   `yaml:"podStrategy"`       // Default: LocalAndRemote
}

type LiqoOpt struct {
	Consumer     string        `yaml:"consumer"`
	Offload      []LiqoOffload `yaml:"offload,omitempty"`
	RuntimeClass bool          `yaml:"runtimeClass"` // Default: empty
}

// Node struct contains generic nodes information
type Node struct {
	Filename string   `yaml:"filename"`
	Count    int      `yaml:"count"`
	Args     []string `yaml:"args,omitempty"`
	name     string   `yaml:"-"`
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
	nodeName     string
	currentIndex int
}

// GetName returns the name from metadata of the
// corresponding node
func (node *Node) GetName() (string, error) {
	if node.name == "" {
		return node.inferName()
	} else {
		return node.name, nil
	}
}

// inferName gets the name of the node for the first time
func (node *Node) inferName() (string, error) {
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
	node.name = ni.Metadata.Name
	return ni.Metadata.Name, nil
}

// GetCurrentIndex fetches the current nodes number from
// the associated node
func (node *Node) GetCurrentIndex(clusterIndex int) int {
	nodeName, _ := node.GetName()
	for i := range nodeCurrentReplicasVec[clusterIndex] {
		if nodeCurrentReplicasVec[clusterIndex][i].nodeName == nodeName {
			return nodeCurrentReplicasVec[clusterIndex][i].currentIndex
		}
	}
	node.SetCurrentIndex(0, clusterIndex)
	return 0
}

// SetCurrentIndex sets the number of nodes deployed from the
// associated node
func (node *Node) SetCurrentIndex(index int, clusterIndex int) {
	nodeName, _ := node.GetName()
	for i := range nodeCurrentReplicasVec[clusterIndex] {
		if nodeCurrentReplicasVec[clusterIndex][i].nodeName == nodeName {
			nodeCurrentReplicasVec[clusterIndex][i].currentIndex = index
			return
		}
	}
	nodeCurrentReplicasVec[clusterIndex] = append(nodeCurrentReplicasVec[clusterIndex], nodeCurrentReplicas{
		nodeName:     nodeName,
		currentIndex: 0,
	})
}

// GetConfName returns the config file name of the
// associated node
func (node *Node) GetConfName() string {
	return node.Filename
}

// GetCount returns the replicas wanted for the associated node
func (node *Node) GetCount() int {
	return node.Count
}

// GetClusterNames returns the cluster's name
func GetClusterNames() []string {
	clusterName := make([]string, len(conf.Clusters))
	for i := range conf.Clusters {
		clusterName[i] = conf.Clusters[i].ClusterName
	}
	return clusterName
}

func GetClusterName(clusterIndex int) string {
	if clusterIndex >= len(conf.Clusters) {
		return ""
	}
	return conf.Clusters[clusterIndex].ClusterName
}

func GetLiqoConsumerCluster() (string, int) {
	return conf.Liqo.Consumer, conf.LiqoConsumerIndex
}

func checkLiqoConsumer() bool {
	for i, cluster := range GetClusterNames() {
		if cluster == conf.Liqo.Consumer {
			conf.LiqoConsumerIndex = i
			return true
		}
	}
	return false
}

func GetLiqoConf() LiqoOpt {
	return conf.Liqo
}

func IsLiqoActive() bool {
	return conf.Liqo.Consumer != ""
}

// GetKwokConf returns the kwok config file name
func GetKwokConf() [][]string {
	kwokConf := make([][]string, len(commandsConf))
	for i := range commandsConf {
		kwokConf[i] = conf.Clusters[i].KwokConfigs
	}
	return kwokConf
}

// GetNodesConf returns a list of Node arranged per cluster
func GetNodesConf() [][]Node {
	var nodes [][]Node
	for i := range commandsConf {
		nodes = append(nodes, conf.Clusters[i].Nodes)
	}
	return nodes
}

// GetAuditConf returns the audit config file name
func GetAuditConf() []string {
	audits := make([]string, len(conf.Clusters))
	for i := range conf.Clusters {
		audits[i] = conf.Clusters[i].Audit
	}
	return audits
}

// GetClusterNameKubeconf returns the cluster's name in the kubeconfig
func GetClusterNameKubeconf() []string {
	res := make([]string, len(conf.Clusters))
	for i := range conf.Clusters {
		res[i] = ClusterKubeconfigs[i].Clusters[0].Name
	}
	return res
}
