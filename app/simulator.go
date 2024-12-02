package app

import (
	"log"
	"main/configs"
	"main/utils"
	"os"
	"path"
)

// Simulation function
//
// starts the simulation by creating the fake cluster,
// running all the commands and deleting it at last, logging everything
// in the logs folder
func Simulation() {
	configs.NewConfig()
	home, _ := os.UserHomeDir()

	// checks if kwokctl and kubectl is installed
	if !utils.CommandExists("kwokctl") {
		log.Fatal("kwokctl not installed")
	}
	if !utils.CommandExists("kubectl") {
		log.Fatal("kubectl not installed")
	}

	go utils.BufferOutWriter()
	go utils.BufferErrWriter()

	// Cluster Creation
	utils.KwokctlCreate()

	// node Creation
	utils.NodeCreate(configs.GetNodesConf())

	// Executes the commands with the specified delay
	utils.ConcurrentQueueRun(configs.GetQueues())

	// Copy and compress log file
	utils.Compress("audit.log", path.Join(home, ".kwok/clusters", configs.GetClusterName(), "logs"))

	// Cluster Deletion
	utils.KwokctlDelete()
}
