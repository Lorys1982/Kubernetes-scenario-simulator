package global

import (
	"time"
)

var LogTime = time.Now().Format("2006-01-02_15:04:05")
var ConfName []string
var StartTime time.Time
var ClusterNames []string

type LogCommandInfo struct {
	CmdIndex     int
	QueueName    string
	ClusterIndex int
}
