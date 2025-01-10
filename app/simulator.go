package app

import (
	"fmt"
	"log"
	"main/apis/v1alpha1"
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
	v1alpha1.NewConfig()
	global.ConfName = v1alpha1.GetCommandsConfName()
	global.ClusterNames = v1alpha1.GetClusterNames()
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
	v1alpha1.ConfPostprocess()

	// If liqo flag is set, install liqo in all clusters, peer consumer and providers
	if v1alpha1.IsLiqoActive() {
		LiqoInstallAll()
		LiqoPeerAll()
		LiqoOffload()
	}

	// node Creation per cluster
	for i := range global.ClusterNames {
		nodes := v1alpha1.GetNodesConf()[i]
		NodeCreate(nodes, i)
	}

	// Executes the commands with the specified delay
	ConcurrentQueueRun(v1alpha1.GetQueues())

	// Cluster Deletion
	KwokctlDeleteAll()
}
