package app

import (
	"gopkg.in/yaml.v2"
	"log"
	"main/configs"
	"os"
)

func Init() {
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

	file, err := os.OpenFile("configs/config.yaml", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0664)

	if !os.IsExist(err) {
		if err != nil {
			log.Fatal(err)
		} else {
			defer file.Close()
			enc := yaml.NewEncoder(file)
			err = enc.Encode(configs.Config{
				ClusterName: "",
				Scheduler:   "",
				Nodes: []configs.Node{
					{
						ConfigName: "example.yaml",
						Count:      0,
					},
				},
				Audit:    "",
				Commands: "",
			})
			if err != nil {
				log.Fatal(err)
			}
			defer enc.Close()
		}
	}

	file, err = os.OpenFile("configs/command_configs/config.yaml", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0664)
	if !os.IsExist(err) {
		if err != nil {
			log.Fatal(err)
		} else {
			defer file.Close()
			enc := yaml.NewEncoder(file)
			err = enc.Encode(configs.Commands{
				Kind:       "",
				ApiVersion: "",
				Metadata: struct {
					Name string `yaml:"name"`
				}{},
				Spec: []configs.CommandsList{
					{
						Exec:  "<command>",
						Delay: 0,
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
