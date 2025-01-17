# Kubernetes Scenario Simulator based on KWOK

This tool is a CLI simulator of K8s **scenarios**, meaning you can plan and configure the **infrastructure** of your cluster based 
on your needs, and then you can simulate the ***execution of commands*** at specific _**times**_ and in different _**queues**_
with each their own _kubeconfig_ file. Finally, the simulator will save **logs** of what happened.

## Project Scope

The simulator can be used in many different testing scenarios, like testing some **custom component** and it's reaction to 
specific scenarios, or even just testing a **scenario itself** before performing it on a real cluster.

Our aim is to provide developers with a **faster** and more **user-friendly** ways to program whatever scenario they have in mind
and to make it reproducible to speed up their developing/testing journey.

## Installation

The simulator is a single executable that can be run on any linux platform, the installation options are:

- **Clone** this repository and build it with **Golang** (```go build``` in the project folder)
- **Download** the pre-compiled binary from releases (_yet to be added_)

## Usage

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
./binary_file -i  # or: ./binary_file --init
```
Enter the `./config` directory and fill the `config.yaml` to define the cluster **topology**:
```bash
cd configs
nano config.yaml # or: whatever other text editor
```

> [!NOTE]
> This is how to fill the config file.
> A more in-depth explanation is below ([Jump](#configuration))
> ```yaml
> apiVersion: k8s-sim.fbk.eu/v1alpha1
> kind: SimConfiguration
> liqo: # Options for Liqo setup
>   consumer: "<consumer-cluster>" # The consumer cluster, if set Liqo will be used
>   offload: # Options for offloading namespaces
>     - namespace: "<namespace-to-offload>" # Namespace to offload
>       clusterSelector: # list of cluster selectors
>         - "<example-selector>"
>       namespaceStrategy: "<ns-strategy>" # Strategy on how to handle namespace names
>       podStrategy: "<pod-strategy>" # Strategy on where to deploy pods
>   runtimeClass: bool # If set, liqo runtimeClass will be created
> clusters: # List of clusters to provision
>   - clusterName: "<cluster-name>" # The name of the cluster
>     kwokConfigs: # The configs to give to kwokctl (for example the configs for a custom scheduler)
>       - "<conf1.yaml>"
>       - "<conf2.yaml>"
>     nodes: # The nodes that you want in the topology
>       - filename: "<node-conf.yaml>" # The config file for the node
>         count: int # the number of nodes you want
>     auditLoggingConfig: "<audit-conf.yaml>" # the audit log config file (same as standard k8s)
>     commandsConfig: "<config.yaml>" # The name of the scenario simulation config file, it is config.yaml by default
> ```
> Make sure to put all the configs written in this config file inside the `./configs/topology` **directory**

Enter the `./configs/command_configs` directory and fill the `config.yaml` to program the **scenario** to simulate:
```bash
cd command_configs
nano config.yaml # or: whatever other text editor
```

> [!NOTE] 
> This is how to fill the config file
> A more in-depth explanation is below ([Jump](#configuration))
> ```yaml
> apiVersion: k8s-sim.fbk.eu/v1alpha1
> kind: SimCommandsConfiguration
> metadata:
> name: "" # The name of the scenario simulation
> spec:
>   queues: # Contain each queue that will run in parallel
>     - name: "" # The name of the queue
>       kubeconfig: "" # The kubecofnig you want to use for this queue (if empty will default to the cluster's ~/.kwok/clusters/<cluster-name>/kubeconfig.yaml)
>       sequence: # The sequence of commands. 
>         # Use exec and command keywords separately
>         - exec: "<command>" # Exec will execute the complete command you provide with no interference from the simulator
>           time: 0 # The absolute time (from the start of the simulation) to which run the command
>           context: "<context-name>" # [Optional] The name of the context in which to run the command
>         - command: "<resource action>" # Command is a wrapper command
>           time: 0
>           filename: "<filename>" # If the command requires a file, put it here
>           resource: "<resource>" # If the command operates on a resource, put it here
>           count: 1 # If the command requires a count put it here
>           args: # If the command requires arguments put them here
>             - "<arg1>"
>           context: "<context-name>" # Context in which execute the command, defaults to the first in the kubeconfig
>           namespace: "<namespace-name> " # Namespace in which execute the command
> ```
> Make sure to put every config written in this file inside the `./configs/command_configs` directory

**Run** the simulation
```bash
./binary_file -s # or: ./binary_file --start
```

Get the **logs** by entering the `./logs` directory, here there will be a directory for each created
cluster, called as said `cluster`.  
Each of these directory will contain the following log files:

- **StdOut Logs** 
- **StdErr Logs**
- **Audit Logs** (compressed .gz)

All three tied to the **execution time** and the **name** of the scenario

## Liqo integration

It is possible to use Liqo in the simulation with the options in a filed called `liqo`, inside the
`SimCongiguration` config file.  
Only a single **consumer** cluster can be specified, which will **peer** to every other cluster of the simulation, and will
then **offload** the specified **namespace** to the selected clusters.

> [!NOTE]
> When using Liqo, the cluster will be set up with Kind to host all the required Liqo components,
> and will then use in-cluster Kwok for every subsequent operations. This means that the startup
> of the cluster will be significantly slower than a normal simulation.

## Configuration

Besides the normal _KWOK_ and _Kubectl_ config files, There are **Two** required config files specific for the simulator.

They will be automatically generated on [**_Initialization,_**](#run-on-linux) and you'll only need to fill them.

### SimConfiguration

**Just one** of this configs file can exist at a time.  
This configuration manages the cluster's topology, it is used to set up one or more Kwok clusters specifying **nodes**,
**cluster configs**, like audit policies and custom components, to set up Liqo if needed, specifying **consumer**, and
**offloading** options, while also containing the name of the **Scenario configs.**

**Location:** `./configs/config.yaml`

#### Fields
- **apiVersion [string]:** Version of the Simulator api
- **kind [string]:** The type of object you are specifying (`SimConfiguration`)
- **liqo [liqo options]:** This field contains the options about Liqo set up 
  - **consumer [string]:** The consumer cluster, must be one of the clusters in the `clusters` field
  - **offload [option list]:** List of options about Liqo offloading
    - **namespace [string]:** The name of the namespace to offload, will be created if not exist (**required**)
    - **clusterSelector [Optional] [string list]:** List of cluster selectors to include/exclude clusters in the offloading process (defaults to `empty`)
    - **namespaceStrategy [Optional] [string]:** Strategy to define the name of the offloaded namespace in the provider clusters,
    the values are: `DefaultName`(default), `EnforceSameName` and `SelectedName`
    - **podStrategy [Optinal] [string]:** Strategy on pod scheduling, the values are: `LocalAndRemote`(default),
    `Local` and `Remote`
  - **runtimeClass [Optional] [bool]:** If set, the runtimeClass flag will be set and the `liqo` runtimeClass will be created
- **clusters [cluster list]:** List of clusters, each with its own `SimCommandsConfiguration`
  - **clusterName [string]:** Name of the kwok cluster, mainly introduced to match the name inside kwok configs
  - **kwokConfigs [string list]:** List of config files (.yaml) applicable to kwok cluster creation 
  - **nodes [nodes list]:** List of nodes, each composed of _name_ and _count_
    - **filename [string]:** Configuration file name of the node to deploy
    - **count [int]:** How many nodes to replicate
  - **auditLoggingPolicy [string]:** Config file of the audit policy (it's the same as standard k8s)
  - **commandsConfigs [string]:** Scenario config file, defaults to `config.yaml` but it's customizable if you want
  to create multiple scenario configs, as long as they are structured correctly and in the `./configs/command_configs` 
  directory

> [!NOTE]
> Nodes will be explained better below [(Jump)](#nodes)

### SimCommandsConfiguration

**Multiple** of these config files can exist at a time. 
This configuration manages the scenario you want to reproduce, it supports multiple **simultaneous queues** to simulate
different users, each with its own **kubeconfig** and **commands sequence** which will execute at a given **time**
and in a specified cluster **context** from the start of the simulation.

**Location:** `./configs/command_configs/config.yaml` (Note that the config name is **variable**)

#### Fields
- **apiVersion [string]:** Version of the Simulator api
- **kind [string]:** The type of object you are specifying (`SimCommandsConfiguration`)
- **metadata [metadata list]:** Data about the scenario
  - **name [string]:** Name of the scenario, will influence the name of the logs [(Jump)](#logs)
- **spec [specs]:** Configs to program the simulation
  - **aliases [string list]:** #TODO (Sets user defined aliases to use with the command field)
  - **queues [queue list]:** List of queues which will run simultaneously  
    - **name [string]:** Name of the queue, will appear inside the logs [(Jump)](#logs)
    - **kubeconfig [string]:** Custom kubeconfig file, if left empty it will take the standard `~/.kube/config`
    - **sequence [commands list]:** Sequence of commands to run, the formats of the commands are the 2 below and cannot be mixed  
      (Note that each command will be run in the `./configs/command_configs` directory)
      - **Exec Format [RAW commands]:** RAW shell commands (Described below)
      - **Command Format [wrapped commands]:** Wrapped commands (Described below)

#### Exec Format
- **exec [string]:** Full RAW shell command
- **time [float]:** Absolute time of execution after the start of the simulation in seconds
- **context [Optional] [string]:** The context into which run the command, if not specified the command will be run
    in the first available context of the given `kubeconfig` file, relative to the `SimCommandsConfig` cluster.
- **namespace [Optional] [string]:** The namespace into which run the command

#### Command Format
- **command [string]:** Simulator provided command, will be explained below [(Jump)](#commands)
- **time [float]:** Absolute time of execution after the start of the simulation in seconds
- **filename [string]:** If the command requires a file, it can be written here, and it will be taken form the 
`commands_conf` directory (i.e. `kube apply`)
- **resource [string]:** If the command requires a resource, it can be written here (i.e. `kube get`)
- **count [int]:** If the command requires a count, it can be written here (i.e. `node create`)
- **args [string list]:** If the command requires arguments, they can be listed here(i.e. `namespace create`)
- **context [Optional] [string]:** The context into which run the command, if not specified the command will be run
in the first available context of the given `kubeconfig` file, relative to the `SimCommandsConfig` cluster.
- **namespace [Optional] [string]:** The namespace into which run the command

### Nodes

Nodes have their own format in the configuration, thanks to that you can **replicate** multiple nodes from a single
node config file.  
The nodes created this way will be called **"_node_name_-N"** with **N** being an incremental number.

> [!WARNING]
> It is highly advised to manage nodes **ONLY** through our **interface** in the command sequence, since if you interfere
> with the incremental numbers modifying in any way the intermediate nodes, the simulator will **not** be able
> to manage nodes automatically anymore, and it will most probably **break**.

### Commands

The simulator provides commands of its own to improve **QoL**.  
These commands are composed by two parts (Order matters): 

`command: <object> <action>`

**Objects:** The _first part_ of the command, determines on what we are acting
  - _**Node** object:_ **Actions** will be executed on **node** resources
  - _**Kube** object:_ **Actions** will be executed using **kubectl**
  - _**Namespace**_ object: **Actions** will be executed on **namespace** resources

**Action:** The _second part_ of the command, determines what to do
  - _**Create** action:_ **Creates** The resource / By using the resource.
  - _**Apply** action:_ **Applies** The resource / By using the resource.
  - _**Delete** action:_ **Deletes** The resource / By using the resource.
  - _**Get** action:_ **Gets** The resource / By using the resource.
  - _**Scale** action:_ **Scales** The resource / by using the resource.

Besides the fields that both `Exec` and `Command` share (**time**, **context** and **namespace**), the other fields are used
differently depending on the **object** and the **action**.

They are summarized in the following table:

|       \       | Create                                                                                                  | Apply                                                                                | Delete                                                                                             | Get                                                    | Scale                                                                                                                                  |
|:-------------:|:--------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------|--------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------|
|   **Node**    | `Filename`: Node to create<br>`Count`: Number of creations<br>`Args`: Additional args                   | `Filename`: Node to apply<br>`Count`: Number of creations<br>`Args`: Additional args | `Filename`: Node to delete<br>`Count`: Number of deletion<br>`Args`: Additional args               | `Args`: Additional args                                | `Filename`: Node to scale<br>`Count`: Target number of replicas<br>`Args`: Additional args                                             |
|   **Kube**    | `Filename`: YAML file to create<br>`Resource`: Resource to create<br>`Args`: Required + additional args | `Filename`: YAML to apply<br>`Args`: Additional args                                 | `Filename`: YAML to delete<br>`Resource`: Resource to delete<br>`Args`: Required + additional args | `Resource`: Resource to get<br>`Args`: Additional args | `Filename`: YAML to scale<br>`Resource`: Resource to scale<br>`Count`: Target number of replicas<br>`Args`: Required + additional args |
| **Namespace** | `Filename`: Namespace file<br>`Args`: Required + additional args                                        | **NONE**                                                                             | `Filename`: Namespace file<br>`Args`: Required + additional args                                   | `Args`: Additional args                                | **NONE**                                                                                                                               |

## Logs

Logs are created to be **easily navigated**, each time a command is executed a prefix will be prepended with the following form:

```text
[Queue: <Q>][Command #<N> Start] 20:57:12.113677 <Command Executed>
[Queue: <Q>][Command #<N>] <StdOut/StdErr>
[Queue: <Q>][Command #<N> End] 20:57:12.769539 Executed at Time: <Time of execution> Seconds
```

- **Q** can be the **name** of the executed **queue**, or a tailor made name, if the command is part of the topology management
- **N** can be the **position** of the command in the sequence, or **0** if the commands is part of the topology management

**Start** is when the command is executed and **End** is when it returns, everything **In Between** is whatever the 
command returned (Even multilined).
