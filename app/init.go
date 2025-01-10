package app

import (
	"gopkg.in/yaml.v2"
	"log"
	"main/configs"
	"os"
)

// Init Function
//
// Initialized config files and directories
func Init() {
	// Directories Creation
	err := os.MkdirAll("./logs", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll("./configs/topology", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll("./configs/command_configs", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// Main config file template creation
	file, err := os.OpenFile("configs/config.yaml", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0664)
	if !os.IsExist(err) {
		if err != nil {
			log.Fatal(err)
		} else {
			defer file.Close()
			enc := yaml.NewEncoder(file)
			err = enc.Encode(configs.Config{
				Kind:       "SimConfiguration",
				ApiVersion: "k8s-sim.fbk.eu/v1alpha1",
				Liqo: configs.LiqoOpt{
					Consumer: "Cluster1",
					Offload: []configs.LiqoOffload{
						{
							Namespace:         "default",
							ClusterSelector:   []string{"selector1", "selector2"},
							NamespaceStrategy: "DefaultName",
							PodStrategy:       "LocalAndRemote",
						},
					},
					RuntimeClass: true,
				},
				Clusters: []configs.Cluster{
					{
						ClusterName: "Cluster1",
						KwokConfigs: []string{
							"--config exampleConf.yaml",
						},
						Nodes: []configs.Node{
							{
								Filename: "example.yaml",
								Count:    0,
							},
						},
						Audit:    "",
						Commands: "config.yaml",
					},
					{
						ClusterName: "Cluster2",
						KwokConfigs: []string{
							"",
						},
						Nodes: []configs.Node{
							{
								Filename: "example.yaml",
								Count:    0,
							},
						},
						Audit:    "",
						Commands: "config.yaml",
					},
				},
			})
			if err != nil {
				log.Fatal(err)
			}
			defer enc.Close()
		}
	}

	// CommandsConf config file template creation
	file, err = os.OpenFile("configs/command_configs/config.yaml", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0664)
	if !os.IsExist(err) {
		if err != nil {
			log.Fatal(err)
		} else {
			defer file.Close()
			enc := yaml.NewEncoder(file)
			err = enc.Encode(configs.CommandsConf{
				Kind:       "SimCommandsConfiguration",
				ApiVersion: "k8s-sim.fbk.eu/v1alpha1",
				Metadata: struct {
					Name string `yaml:"name"`
				}{},
				Spec: struct {
					Aliases []string        `yaml:"aliases,omitempty"`
					Queues  []configs.Queue `yaml:"queues"`
				}{
					Queues: []configs.Queue{
						{
							Name:       "",
							Kubeconfig: "",
							Sequence: []configs.Command{
								{
									Exec:      "command",
									Time:      0,
									Context:   "context-name",
									Namespace: "namespace-name",
								},
								{
									Command:  "wrapper command",
									Time:     0,
									Filename: "filename",
									Resource: "resource",
									Count:    1,
									Args: []string{
										"args1",
									},
									Context:   "context-name",
									Namespace: "namespace-name",
								},
							},
						},
					},
				},
			})
			if err != nil {
				log.Fatal(err)
			}
			defer enc.Close()
		}
	}
}
