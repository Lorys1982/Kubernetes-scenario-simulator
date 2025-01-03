package app

import (
	"main/configs"
	"os/exec"
	"sync"
)

func LiqoInstallAll() {
	clusters := configs.GetClusterNameKubeconf()
	wg := &sync.WaitGroup{}
	for i := range clusters {
		wg.Add(1)
		go liqoInstall(i, wg)
	}
	wg.Wait()
}

func liqoInstall(clusterIndex int, wg *sync.WaitGroup) {
	info := commandInfo{
		Queue: configs.Queue{
			Name:        "LiqoInstall",
			Kubeconfig:  configs.GetKubeConfigPath(clusterIndex),
			KubeContext: configs.ClusterKubeconfigs[clusterIndex].Contexts[0],
		},
		CmdIndex: 0,
	}
	nodeName := "kwok" + "-" + configs.GetClusterName(clusterIndex) + "-" + "control-plane"
	KubectlUncordon(0, info, nodeName)
	cmd := exec.Command("liqoctl", "install", "kind", "--context", info.Queue.KubeContext.Name, "--disable-kernel-version-check")
	err := commandRun(cmd, 0, info)
	if err != nil {
		crashLog(err.Error(), info)
	}
	wg.Done()
}

func LiqoPeerAll() {
	consumerCluster, consumerIndex := configs.GetLiqoConsumerCluster()
	wg := &sync.WaitGroup{}
	info := commandInfo{
		Queue: configs.Queue{
			Name:        "LiqoPeer",
			Kubeconfig:  configs.GetKubeConfigPath(consumerIndex),
			KubeContext: configs.ClusterKubeconfigs[consumerIndex].Contexts[0],
		},
		CmdIndex: 0,
	}
	for _, cluster := range configs.GetClusterNames() {
		if cluster != consumerCluster {
			wg.Add(1)
			go liqoPeer(cluster, info, wg)
		}
	}
	wg.Wait()
}

func liqoPeer(cluster string, info commandInfo, wg *sync.WaitGroup) {
	cmd := exec.Command("liqoctl", "peer", "--remote-kubeconfig", configs.GenKubeConfigPath(cluster),
		"--server-service-type", "NodePort", "--context", info.Queue.KubeContext.Name)
	err := commandRun(cmd, 0, info)
	if err != nil {
		crashLog(err.Error(), info)
	}
	wg.Done()
}

func LiqoOffload() {
	_, consumerIndex := configs.GetLiqoConsumerCluster()
	var args []string
	info := commandInfo{
		Queue: configs.Queue{
			Name:        "LiqoOffload",
			Kubeconfig:  configs.GetKubeConfigPath(consumerIndex),
			KubeContext: configs.ClusterKubeconfigs[consumerIndex].Contexts[0],
		},
		CmdIndex: 0,
	}
	for _, offloadInfo := range configs.GetLiqoConf().Offload {
		args = append(args, "offload", "namespace", offloadInfo.Namespace)
		if len(offloadInfo.NamespaceStrategy) != 0 {
			args = append(args, "--namespace-mapping-strategy", offloadInfo.NamespaceStrategy)
		}
		for i := range offloadInfo.ClusterSelector {
			args = append(args, "--selector", offloadInfo.ClusterSelector[i])
		}
		if len(offloadInfo.PodStrategy) != 0 {
			args = append(args, "--pod-offloading-strategy", offloadInfo.PodStrategy)
		}
		args = append(args, "--context", info.Queue.KubeContext.Name)
		cmd := exec.Command("liqoctl", args...)
		err := commandRun(cmd, 0, info)
		if err != nil {
			crashLog(err.Error(), info)
		}
	}
	for clusterIndex := range configs.GetClusterNames() {
		nodeName := "kwok" + "-" + configs.GetClusterName(clusterIndex) + "-" + "control-plane"
		info = commandInfo{
			Queue: configs.Queue{
				Name:        "LiqoEnd",
				Kubeconfig:  configs.GetKubeConfigPath(clusterIndex),
				KubeContext: configs.ClusterKubeconfigs[clusterIndex].Contexts[0],
			},
			CmdIndex: 0,
		}
		KubectlCordon(0, info, nodeName)
	}
}
