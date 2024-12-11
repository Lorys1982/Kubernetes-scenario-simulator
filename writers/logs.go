package writers

import (
	"bytes"
	"fmt"
	"github.com/notEpsilon/go-pair"
	"log"
	"main/global"
	"os"
)

var LogChannelStd = make(chan pair.Pair[[]byte, int])
var LogChannelErr = make(chan pair.Pair[[]byte, int])
var killChannelErr = make(chan bool)

func BufferOutWriter() {
	var toWrite []byte
	var clusterIndex int
	var outFiles []*os.File
	for i, cluster := range global.ClusterNames {
		outFiles[i], _ = os.OpenFile(fmt.Sprintf("logs/%s/%s_StdOut_%s.log", cluster, global.ConfName, global.LogTime), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	}
	for {
		select {
		case toWrite, clusterIndex = <-LogChannelStd:
			outFiles[clusterIndex].Write(toWrite)
		}
	}
}

func BufferErrWriter() {
	var toWrite []byte
	var clusterIndex int
	var errFiles []*os.File
	for i, cluster := range global.ClusterNames {
		errFiles[i], _ = os.OpenFile(fmt.Sprintf("logs/%s/%s_StdErr_%s.log", cluster, global.ConfName, global.LogTime), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	}
	for {
		select {
		case <-killChannelErr:
			return
		case toWrite, clusterIndex = <-LogChannelErr:
			errFiles[clusterIndex].Write(toWrite)
		}
	}
}

func CrashLog(err string, clusterIndex int, option ...func()) {
	var buffer bytes.Buffer
	errLog := log.New(&buffer, "[Fatal Error] ", log.Ltime|log.Lmicroseconds)
	errLog.Printf(err + "\n\n")

	if global.ConfName == nil {
		errFile, _ := os.OpenFile(fmt.Sprintf("logs/%s_StdErr_%s.log", global.ConfName, global.LogTime), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		errFile.Write(buffer.Bytes())
		log.Fatal(buffer.String())
	} else {
		LogChannelErr <- *pair.New(buffer.Bytes(), clusterIndex)
	}

	if option != nil {
		f := option[0]
		f()
	}
	killChannelErr <- true
	log.Fatal(buffer.String())
}

func ErrLog(err string, cmd string, clusterIndex int) {
	var buffer bytes.Buffer
	errLog := log.New(&buffer, "[Error] ", log.Ltime|log.Lmicroseconds)
	errLog.Printf("(Command: %s) %s\n\n", cmd, err)
	LogChannelErr <- *pair.New(buffer.Bytes(), clusterIndex)
}
