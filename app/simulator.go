package app

import (
	"log"
	"main/configs"
	"main/global"
	"main/writers"
)

// Simulation function
//
// starts the simulation by creating the fake cluster,
// running all the commands and deleting it at last, logging everything
// in the logs folder
func Simulation() {
	configs.NewConfig()
	global.ConfName = configs.GetCommandsConfName()

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

	// Cluster Deletion
	KwokctlDelete()
}
