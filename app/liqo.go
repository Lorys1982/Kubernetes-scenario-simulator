package app

import (
	"main/configs"
	"os/exec"
	"sync"
)

func LiqoInstallAll() {
	clusters := configs.GetClusterNameKubeconf()
	wg := &sync.WaitGroup{}
	for i, cluster := range clusters {
		wg.Add(1)
		go liqoInstall(cluster, i, wg)
	}
	wg.Wait()
}

func liqoInstall(cluster configs.KubeCluster, clusterIndex int, wg *sync.WaitGroup) {
	info := commandInfo{
		Queue: configs.Queue{
			Name:       "Liqo",
			Kubeconfig: "",
			KubeContext: configs.Context{
				ClusterIndex: clusterIndex,
			},
		},
		CmdIndex: 0,
	}
	nodeName := "kwok" + "-" + configs.GetClusterName(clusterIndex) + "-" + "control-plane"
	KubectlUncordon(0, commandInfo{
		Queue: configs.Queue{
			Name:        "Liqo",
			KubeContext: configs.DefaultKubeconfig.Contexts[clusterIndex],
		},
		CmdIndex: 0,
	}, nodeName)
	cmd := exec.Command("liqoctl", "install", "kind", "--cluster", cluster.Name, "--disable-kernel-version-check")
	err := commandRun(cmd, 0, info)
	if err != nil {
		crashLog(err.Error(), info)
	}
	wg.Done()
}

func LiqoPeer() {

}

func LiqoOffload() {

}
