package main

import (
	"fmt"
	_ "k8s.io/client-go"
	"log"
	"main/Utils"
	"os/exec"
	"time"
)

func main() {
	command := "kwokctl"

	// checks if kwokctl is installed
	if !Utils.CommandExists(command) {
		log.Fatal(command, " not installed")
	}

	cmd := exec.Command("kwokctl", "create", "cluster")
	Utils.CommandRun(cmd)

	for v := 0; v < 10; v++ {
		fmt.Print(v, " ")
		time.Sleep(1 * time.Second)
	}

	cmd = exec.Command("kwokctl", "delete", "cluster")
	Utils.CommandCleanRun(cmd)
}
