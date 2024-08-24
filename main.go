package main

import (
	"fmt"
	"main/app"
	"os"
)

func main() {

	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "-h" || argsWithoutProg[0] == "--help" {
		fmt.Println("Usage: 	./<binary>")
		fmt.Println(" -h	--help show this output")
		fmt.Println(" -i	--init the environment")
		fmt.Println(" -s	--start the simulation")
		fmt.Println(" -c	--commands available as wrappers")
		fmt.Println("Note that only one argument will be taken into consideration (the first)")
	} else if argsWithoutProg[0] == "--init" || argsWithoutProg[0] == "-i" {
		app.Init()
	} else if argsWithoutProg[0] == "--start" || argsWithoutProg[0] == "-s" {
		app.Simulation()
	} else if argsWithoutProg[0] == "--commands" || argsWithoutProg[0] == "-c" {
		fmt.Println("Commands:")
		fmt.Println(" Create		Creates the resource")
		fmt.Println("   Node		Node resource")
		fmt.Println("     filename	The name of the node .yaml file")
		fmt.Println("     count	How many nodes to create")
		fmt.Println(" Delete		Deletes the resource")
		fmt.Println("   Node		Node resource")
		fmt.Println("     filename	The name of the node .yaml file")
		fmt.Println("     count	How many nodes to delete")
		fmt.Println("Example: \n" +
			"- command: create node\n" +
			"  filename: node.yaml\n" +
			"  count: 2")
		fmt.Println("Note that the order of the command is irrelevant (create node = node create)")
	} else {
		fmt.Println("Unknown argument: " + argsWithoutProg[0])
	}
}
