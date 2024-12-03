package app

import (
	"log"
	"main/configs"
	"main/constants"
	"main/utils"
	"main/writers"
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
	constants.ConfName = configs.GetCommandsConfName()
	home, _ := os.UserHomeDir()

	// checks if kwokctl and kubectl is installed
	if !CommandExists("kwokctl") {
		log.Fatal("kwokctl not installed")
	}
	if !CommandExists("kubectl") {
		log.Fatal("kubectl not installed")
	}

	// Logger Initialization
	go writers.BufferOutWriter()
	go writers.BufferErrWriter()

	// Cluster Creation
	KwokctlCreate()

	// node Creation
	NodeCreate(configs.GetNodesConf())

	// Executes the commands with the specified delay
	ConcurrentQueueRun(configs.GetQueues())

	// Copy and compress log file
	utils.Compress("audit.log", path.Join(home, ".kwok/clusters", configs.GetClusterName(), "logs"))

	// Cluster Deletion
	KwokctlDelete()
}
