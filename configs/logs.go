package configs

import (
	"fmt"
	"log"
	"os"
)

func CrashLog(err string) {
	errFile, _ := os.OpenFile(fmt.Sprintf("logs/%s_StdErr_%s.log", GetCommandsConfName(), LogTime), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	errLog := log.New(errFile, "[Fatal Error] ", log.Ltime|log.Lmicroseconds)
	errLog.Println(err)
	log.Fatal(err)
}

func ErrLog(err string, cmd string) {
	errFile, _ := os.OpenFile(fmt.Sprintf("logs/%s_StdErr_%s.log", GetCommandsConfName(), LogTime), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	errLog := log.New(errFile, "[Error] ", log.Ltime|log.Lmicroseconds)
	errLog.Printf("(Command: %s) %s\n\n", cmd, err)
}
