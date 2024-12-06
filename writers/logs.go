package writers

import (
	"bytes"
	"fmt"
	"log"
	"main/global"
	"os"
)

var LogChannelStd = make(chan []byte)
var LogChannelErr = make(chan []byte)
var killChannelErr = make(chan bool)

func BufferOutWriter() {
	var toWrite []byte
	outFile, _ := os.OpenFile(fmt.Sprintf("logs/%s_StdOut_%s.log", global.ConfName, global.LogTime), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	defer outFile.Close()
	for {
		select {
		case toWrite = <-LogChannelStd:
			outFile.Write(toWrite)
		}
	}
}

func BufferErrWriter() {
	var toWrite []byte
	errFile, _ := os.OpenFile(fmt.Sprintf("logs/%s_StdErr_%s.log", global.ConfName, global.LogTime), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	defer errFile.Close()
	for {
		select {
		case <-killChannelErr:
			return
		case toWrite = <-LogChannelErr:
			errFile.Write(toWrite)
		}
	}
}

func CrashLog(err string, option ...func()) {
	var buffer bytes.Buffer
	errLog := log.New(&buffer, "[Fatal Error] ", log.Ltime|log.Lmicroseconds)
	errLog.Printf(err + "\n\n")

	if global.ConfName == nil {
		errFile, _ := os.OpenFile(fmt.Sprintf("logs/%s_StdErr_%s.log", global.ConfName, global.LogTime), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		errFile.Write(buffer.Bytes())
		log.Fatal(buffer.String())
	} else {
		LogChannelErr <- buffer.Bytes()
	}

	if option != nil {
		f := option[0]
		f()
	}
	killChannelErr <- true
	log.Fatal(buffer.String())
}

func ErrLog(err string, cmd string) {
	var buffer bytes.Buffer
	errLog := log.New(&buffer, "[Error] ", log.Ltime|log.Lmicroseconds)
	errLog.Printf("(Command: %s) %s\n\n", cmd, err)
	LogChannelErr <- buffer.Bytes()
}
