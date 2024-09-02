# Kubernetes Scenario Simulator based on KWOK

This tool is a CLI simulator of K8s **scenarios**, meaning you can plan and configure the **infrastructure** of your cluster based 
on your needs, and then you can simulate the ***execution of commands*** at specific _**times**_ and in different _**queues**_
with each their own _kubeconfig_ file.  
Finally the simulator will save **logs** of what happened.

## Project Scope

The simulator can be used in many different testing scenarios, like testing some **custom component** and it's reaction to 
specific scenarios, or even just testing a **scenario itself** before performing it on a real cluster.

Our aim is to provide developers with a **faster** and more **user-friendly** ways to program whatever scenario they have in mind
and to make it reproducible to speed up their developing/testing journey.

## Installation

The simulator is a single executable that can be run on any linux platform, the installation options are:

- **Clone** this repository and build it with **Golang** (```go build``` in the project folder)
- **Download** the pre-compiled binary from releases (_yet to be added_)

## How to use

The simulator is entirely **CLI**, but you'll have to use _**yaml**_ configurations to set the scenarios.

### Prerequisites

This simulator is based on the usage of ```KWOK``` and ```Kubectl```, therefore it requires:

- **KWOK** and **Kwokctl** 
- **Kubectl** 
- **Docker** 

> [!WARNING]
> Make sure that all of the above can be executed by your user

### Run on Linux

**Initialize** the required files and directories:
```bash
binary_file -i  # or: binary_file --init
```
Enter the `./config` directory and fill the `config.yaml` to define the cluster topology:
```bash
cd configs
nano config.yaml # or: whatever other text editor
```

> [!NOTE]
> This is how to fill the config file
> ```yaml
> clusterName: "<cluster-name>" # The name of the cluster
> kwokConfigs: # The configs to give to kwokctl (for example the configs for a custom scheduler)
>   - "<conf1.yaml>"
>   - "<conf2.yaml>"
> nodes: # The nodes that you want in the topology
>   - name: "<node-conf.yaml>" # The config file for the node
>     count: int # the number of nodes you want
> auditLoggingConfig: "<audit-conf.yaml>" # the audit log config file (same as standard k8s)
> ```
> Make sure to put all the configs written in this config file inside the `./configs/topology` **directory**

Enter the `./configs/command_configs` directory and fill the `config.yaml` to program the **scenario** to simulate
```bash
cd command_configs
nano config.yaml # or: whatever other text editor
```

> [!NOTE] This is how to fill the config file
> ```yaml
> kind: ""
> apiVersion: ""
> metadata:
> name: ""
> spec:
> aliases: []
> queues:
> - name: ""
>   kubeconfig: ""
>   sequence:
>     - exec: <command>
>       time: 0
>       command: ""
>       filename: ""
>       count: 0
>       args: []
>     - exec: ""
>       time: 0
>       command: <wrapper command>
>       filename: ""
>       count: 0
>       args: []
> ```