package writers

import (
	"bytes"
	"fmt"
	"github.com/notEpsilon/go-pair"
	"log"
	"main/global"
	"main/utils"
	"os"
)

var LogChannelStd = make(chan pair.Pair[[]byte, int])
var LogChannelErr = make(chan pair.Pair[[]byte, int])
var killChannelErr = make(chan bool)

func BufferOutWriter() {
	var pairs pair.Pair[[]byte, int]
	var outFiles = make([]*os.File, len(global.ClusterNames))
	for i, cluster := range global.ClusterNames {
		outFiles[i], _ = os.OpenFile(fmt.Sprintf("logs/%s/%s_StdOut_%s.log", cluster, global.ConfName[i], global.LogTime), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	}
	for {
		select {
		case pairs = <-LogChannelStd:
			outFiles[pairs.Second].Write(pairs.First)
		}
	}
}

func BufferErrWriter() {
	var pairs pair.Pair[[]byte, int]
	var errFiles = make([]*os.File, len(global.ClusterNames))
	for i, cluster := range global.ClusterNames {
		errFiles[i], _ = os.OpenFile(fmt.Sprintf("logs/%s/%s_StdErr_%s.log", cluster, global.ConfName[i], global.LogTime), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	}
	for {
		select {
		case <-killChannelErr:
			return
		case pairs = <-LogChannelErr:
			errFiles[pairs.Second].Write(pairs.First)
		}
	}
}

func CrashLog(err string, info utils.Option[global.LogCommandInfo], option ...func()) {
	var buffer bytes.Buffer

	if global.ConfName == nil {
		errFile, _ := os.OpenFile(fmt.Sprintf("logs/%s_StdErr_%s.log", global.ConfName, global.LogTime), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		errFile.Write(buffer.Bytes())
		log.Fatal(buffer.String())
	} else if info.IsNone() {
		errLog := log.New(&buffer, "[Fatal Error] ", log.Ltime|log.Lmicroseconds)
		errLog.Printf(err + "\n\n")
		for i := range global.ClusterNames {
			LogChannelErr <- *pair.New(buffer.Bytes(), i)
		}
	} else {
		info := info.GetSome()
		errLog := log.New(&buffer, fmt.Sprintf("[Fatal Error][Queue: %s][Command #%d] ", info.QueueName, info.CmdIndex), log.Ltime|log.Lmicroseconds)
		errLog.Printf(err + "\n\n")
		LogChannelErr <- *pair.New(buffer.Bytes(), info.ClusterIndex)
	}

	if option != nil {
		f := option[0]
		f()
	}
	log.Fatal(buffer.String())
}

func ErrLog(err string, cmd string, info global.LogCommandInfo) {
	var buffer bytes.Buffer
	errLog := log.New(&buffer, fmt.Sprintf("[Error][Queue: %s][Command #%d] ", info.QueueName, info.CmdIndex), log.Ltime|log.Lmicroseconds)
	errLog.Printf("(Command: %s) %s\n\n", cmd, err)
	LogChannelErr <- *pair.New(buffer.Bytes(), info.ClusterIndex)
}

func KillLoggers() {
	killChannelErr <- true
}
