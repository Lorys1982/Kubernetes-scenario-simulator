package main

import (
	"fmt"
	_ "k8s.io/client-go"
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

	// TODO whatever operation and instruction i want
	for v := 0; v < 10; v++ {
		fmt.Print(v, " ")
		time.Sleep(1 * time.Second)
	}

	// Cluster Deletion
	utils.KwokctlDelete()

}
