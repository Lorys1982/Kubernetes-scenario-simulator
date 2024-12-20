package app

import (
	"fmt"
	"log"
	"main/configs"
	"main/global"
	"main/writers"
	"os"
)

// Simulation function
//
// starts the simulation by creating the fake cluster,
// running all the commands and deleting it at last, logging everything
// in the logs folder
func Simulation() {
	configs.NewConfig()
	global.ConfName = configs.GetCommandsConfName()
	global.ClusterNames = configs.GetClusterNames()
	for _, cluster := range global.ClusterNames {
		os.MkdirAll(fmt.Sprintf("logs/%s", cluster), os.ModePerm)
	}

	// checks if kwokctl and kubectl is installed
	if !CommandExists("kwokctl") {
		log.Fatal("kwokctl not installed")
	}
	if !CommandExists("kubectl") {
		log.Fatal("kubectl not installed")
	}
	if !CommandExists("liqoctl") {
		log.Fatal("liqoctl not installed")
	}

	// Logger Initialization
	go writers.BufferOutWriter()
	go writers.BufferErrWriter()

	// Cluster Creation
	KwokctlCreateAll()

	// Fill kubeconf structs
	configs.ConfPostprocess()

	// If liqo flag is set, install liqo in all clusters
	if configs.GetLiqoConf() {
		LiqoInstallAll()
	}

	// node Creation per cluster
	for i := range global.ClusterNames {
		nodes := configs.GetNodesConf()[i]
		NodeCreate(nodes, i)
	}

	// Executes the commands with the specified delay
	ConcurrentQueueRun(configs.GetQueues())

	// Cluster Deletion
	KwokctlDeleteAll()
}
