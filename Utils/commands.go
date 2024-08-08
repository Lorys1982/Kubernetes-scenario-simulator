package Utils

import (
	"log"
	"os"
	"os/exec"
)

func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func CommandRun(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func CommandCleanRun(cmd *exec.Cmd) {
	err := cmd.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err.Error())
	}
}
