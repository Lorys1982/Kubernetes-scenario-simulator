package app

import (
	"main/apis/v1alpha1"
	. "main/utils"
	"os/exec"
	"sync"
)

func LiqoInstallAll() {
	clusters := v1alpha1.GetClusterNameKubeconf()
	wg := &sync.WaitGroup{}
	for i := range clusters {
		wg.Add(1)
		go liqoInstall(i, wg)
	}
	wg.Wait()
}

func liqoInstall(clusterIndex int, wg *sync.WaitGroup) {
	info := commandInfo{
		QueueName:   "LiqoInstall",
		Kubeconfig:  v1alpha1.GetKubeConfigPath(clusterIndex),
		KubeContext: v1alpha1.ClusterKubeconfigs[clusterIndex].Contexts[0],
		CmdIndex:    0,
	}
	nodeName := "kwok" + "-" + v1alpha1.GetClusterName(clusterIndex) + "-" + "control-plane"
	args := []string{"install", "kind", "--context", info.KubeContext.Name}
	if v1alpha1.GetLiqoConf().RuntimeClass {
		args = append(args, "--set", "offloading.runtimeClass.enabled=true")
	}
	KubectlUncordon(0, info, nodeName)
	cmd := exec.Command("liqoctl", "install", "kind", "--context", info.KubeContext.Name)
	err := commandRun(cmd, 0, info)
	if err != nil {
		crashLog(err.Error(), info)
	}
	wg.Done()
}

func LiqoPeerAll() {
	consumerCluster, consumerIndex := v1alpha1.GetLiqoConsumerCluster()
	wg := &sync.WaitGroup{}
	info := commandInfo{
		QueueName:   "LiqoPeer",
		Kubeconfig:  v1alpha1.GetKubeConfigPath(consumerIndex),
		KubeContext: v1alpha1.ClusterKubeconfigs[consumerIndex].Contexts[0],
		CmdIndex:    0,
	}
	for _, cluster := range v1alpha1.GetClusterNames() {
		if cluster != consumerCluster {
			wg.Add(1)
			go liqoPeer(cluster, info, wg)
		}
	}
	wg.Wait()
}

func liqoPeer(cluster string, info commandInfo, wg *sync.WaitGroup) {
	cmd := exec.Command("liqoctl", "peer", "--remote-kubeconfig", v1alpha1.GenKubeConfigPath(cluster),
		"--server-service-type", "NodePort", "--context", info.KubeContext.Name)
	err := commandRun(cmd, 0, info)
	if err != nil {
		crashLog(err.Error(), info)
	}
	wg.Done()
}

func LiqoOffload() {
	_, consumerIndex := v1alpha1.GetLiqoConsumerCluster()
	var args []string
	info := commandInfo{
		QueueName:   "LiqoOffload",
		Kubeconfig:  v1alpha1.GetKubeConfigPath(consumerIndex),
		KubeContext: v1alpha1.ClusterKubeconfigs[consumerIndex].Contexts[0],
		CmdIndex:    0,
	}
	for _, offloadInfo := range v1alpha1.GetLiqoConf().Offload {
		KubectlCreate(Some("namespace"), None[string](), 0, info, offloadInfo.Namespace)
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
		args = append(args, "--context", info.KubeContext.Name)
		cmd := exec.Command("liqoctl", args...)
		err := commandRun(cmd, 0, info)
		if err != nil {
			crashLog(err.Error(), info)
		}
	}
	for clusterIndex := range v1alpha1.GetClusterNames() {
		nodeName := "kwok" + "-" + v1alpha1.GetClusterName(clusterIndex) + "-" + "control-plane"
		info = commandInfo{
			QueueName:   "LiqoEnd",
			Kubeconfig:  v1alpha1.GetKubeConfigPath(clusterIndex),
			KubeContext: v1alpha1.ClusterKubeconfigs[clusterIndex].Contexts[0],
			CmdIndex:    0,
		}
		KubectlCordon(0, info, nodeName)
	}
}
