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
	file, err := os.OpenFile("configs/config1.yaml", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0664)
	if !os.IsExist(err) {
		if err != nil {
			log.Fatal(err)
		} else {
			defer file.Close()
			enc := yaml.NewEncoder(file)
			err = enc.Encode(configs.Config{
				Clusters: []configs.Cluster{
					{
						ClusterName: "Cluster1",
						KwokConfigs: []string{
							"",
						},
						Nodes: []configs.Node{
							{
								ConfigName: "example.yaml",
								Count:      0,
							},
						},
						Audit:    "",
						Commands: "config1.yaml",
					},
					{
						ClusterName: "Cluster2",
						KwokConfigs: []string{
							"",
						},
						Nodes: []configs.Node{
							{
								ConfigName: "example.yaml",
								Count:      0,
							},
						},
						Audit:    "",
						Commands: "config1.yaml",
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
	file, err = os.OpenFile("configs/command_configs/config1.yaml", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0664)
	if !os.IsExist(err) {
		if err != nil {
			log.Fatal(err)
		} else {
			defer file.Close()
			enc := yaml.NewEncoder(file)
			err = enc.Encode(configs.CommandsConf{
				Kind:       "",
				ApiVersion: "",
				Metadata: struct {
					Name string `yaml:"name"`
				}{},
				Spec: struct {
					Aliases []string        `yaml:"aliases"`
					Queues  []configs.Queue `yaml:"queues"`
				}{
					Queues: []configs.Queue{
						{
							Name:       "",
							Kubeconfig: "",
							Sequence: []configs.Command{
								{
									Exec:    "<command>",
									Time:    0,
									Context: "<context-name>",
								},
								{
									Command:  "<wrapper command>",
									Filename: "<filename>",
									Count:    1,
									Time:     0,
									Args: []string{
										"<args1>",
									},
									Context: "<context-name>",
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
