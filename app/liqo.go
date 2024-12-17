package app

import (
	"main/configs"
	"os/exec"
)

func LiqoInstall() {
	// TODO IMPORTANT! in some ay this function should respond with the actual kwubeconf names, not hardcoded ones
	clusters := configs.GetClusterNameKubeconf()
	info := commandInfo{
		Queue: configs.Queue{
			Name:       "Liqo",
			Kubeconfig: "",
			KubeContext: configs.Context{
				ClusterIndex: 0,
			},
		},
		CmdIndex: 0,
	}
	for i, cluster := range clusters {
		info.Queue.KubeContext.ClusterIndex = i
		cmd := exec.Command("liqoctl", "install", "--cluster", cluster)
		err := commandRun(cmd, 0, info)
		if err != nil {
			crashLog(err.Error(), info)
		}
	}
}

func LiqoPeer() {

}

func LiqoOffload() {

}
