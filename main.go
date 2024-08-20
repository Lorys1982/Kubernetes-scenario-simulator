package main

import (
	"log"
	"main/configs"
	"main/utils"
	"os"
	"path"
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

	// Executes the commands with the specified delay
	utils.SequentialCommandRun(configs.GetCommandsList())

	// Test multiple node creations
	utils.NodeCreate(configs.GetNodesConf())
	// Test multinode deletion
	utils.NodeDelete(configs.GetNodesConf())

	// Copy and compress log file
	utils.Compress("audit.log", path.Join(home, ".kwok/clusters", configs.GetClusterName(), "logs"))

	// Cluster Deletion
	utils.KwokctlDelete()
}
