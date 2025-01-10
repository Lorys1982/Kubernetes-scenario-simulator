package v1alpha1

const (
	ResourceKindSimCommandsConfiguration = "SimCommandsConfiguration"
)

var commandsConf []*CommandsConf

// CommandsConf struct contains the data of the
// commands configuration file
type CommandsConf struct {
	Kind       string `yaml:"kind"`
	ApiVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Aliases []string `yaml:"aliases,omitempty"`
		Queues  []Queue  `yaml:"queues"`
	}
}

// Command struct contains a single command of the
// command sequence inside a Queue
type Command struct {
	Exec      string   `yaml:"exec,omitempty"`
	Command   string   `yaml:"command,omitempty"`
	Time      float64  `yaml:"time"`
	Resource  string   `yaml:"resource,omitempty"`
	Filename  string   `yaml:"filename,omitempty"`
	Count     int      `yaml:"count,omitempty"`
	Args      []string `yaml:"args,omitempty"`
	Context   string   `yaml:"context,omitempty"`
	Namespace string   `yaml:"namespace,omitempty"`
	index     int      `yaml:"-"`
}

// Queue struct contains a Command sequence and the
// data of the Queue
type Queue struct {
	Name        string    `yaml:"name"`
	Kubeconfig  string    `yaml:"kubeconfig"`
	Sequence    []Command `yaml:"sequence"`
	KubeContext Context   `yaml:"-"`
}

// GetCommandsConfName returns the name of the commands config file
func GetCommandsConfName() []string {
	if commandsConf == nil {
		return make([]string, 0)
	}
	var confNames = make([]string, len(commandsConf))
	for i := range commandsConf {
		confNames[i] = commandsConf[i].Metadata.Name
	}
	return confNames
}

// GetQueues returns a list of Queue per CommandsConf (per cluster)
func GetQueues() [][]Queue {
	var queues [][]Queue
	for i := range commandsConf {
		queues = append(queues, commandsConf[i].Spec.Queues)
	}
	return queues
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
