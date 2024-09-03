package app

import (
	"gopkg.in/yaml.v2"
	"log"
	"main/configs"
	"os"
)

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
				ClusterName: "",
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
				Commands: "config.yaml",
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
									Exec: "<command>",
									Time: 0,
								},
								{
									Command:  "<wrapper command>",
									Filename: "<filename>",
									Count:    1,
									Time:     0,
									Args: []string{
										"<args1>",
									},
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
