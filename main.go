package main

import (
	"fmt"
	"log"
	"main/configs"
	"main/utils"
	"time"
)

func main() {
	command := "kwokctl"
	configs.NewConfig()

	// checks if kwokctl is installed
	if !utils.CommandExists(command) {
		log.Fatal(command, " not installed")
	}

	// Cluster Creation
	utils.KwokctlCreate()

	// Node Creation
	utils.NodeCreate(configs.GetNodesConf())

	// TODO whatever operation and instruction i want
	for v := 0; v < 1; v++ {
		fmt.Print(v, " ")
		time.Sleep(1 * time.Second)
	}

	// Copy and compress log file
	utils.Compress("audit.log")

	// Cluster Deletion
	utils.KwokctlDelete()
}
