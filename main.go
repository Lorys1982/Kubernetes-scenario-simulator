package main

import (
	"fmt"
	"main/app"
	"os"
)

func main() {

	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "-h" || argsWithoutProg[0] == "--help" { // CLI help display
		fmt.Println("Usage: 	./<binary>")
		fmt.Println(" -h	--help show this output")
		fmt.Println(" -i	--init the environment")
		fmt.Println(" -s	--start the simulation")
		fmt.Println("Note that only one argument will be taken into consideration (the first)")
	} else if argsWithoutProg[0] == "--init" || argsWithoutProg[0] == "-i" { // Initialization of conf files and directories
		app.Init()
	} else if argsWithoutProg[0] == "--start" || argsWithoutProg[0] == "-s" { // Simulation start
		app.Simulation()
	} else { // wrong args provided
		fmt.Println("Unknown argument: " + argsWithoutProg[0])
	}
}
