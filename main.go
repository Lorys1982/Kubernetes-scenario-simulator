package main

import (
	"fmt"
	"log"
	"main/configs"
	"main/utils"
	"os"
	"path"
	"time"
)

func main() {
	configs.NewConfig()
	home, _ := os.UserHomeDir()

	// checks if kwokctl and kubectl is installed
	if !utils.CommandExists("kwokctl") {
		log.Fatal("kwokctl not installed")
	}
	if !utils.CommandExists("kubectl") {
		log.Fatal("kubectl not installed")
	}

	// Cluster Creation
	utils.KwokctlCreate()

	// Node Creation
	utils.NodeCreate(configs.GetNodesConf())

	// TODO whatever operation and instruction i want
	for v := 0; v < 1; v++ {
		fmt.Println(v, " ")
		time.Sleep(1 * time.Second)
	}

	// Test multiple node creations
	utils.NodeCreate(configs.GetNodesConf())

	// Copy and compress log file
	utils.Compress("audit.log", path.Join(home, ".kwok/clusters", configs.GetClusterName(), "logs"))

	// Cluster Deletion
	utils.KwokctlDelete()
}
