package main

import (
	"encoding/json"
	"experiment_lwc/commons"
	"experiment_lwc/db"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func HandleListContainer() {
	DB := db.Setup(dbPath)
	tx, err := DB.DB.Begin(true)
	commons.Must(err)
	defer func() {
		_ = tx.Rollback()
		commons.Must(DB.DB.Close())
	}()
	b, err := json.MarshalIndent(DB.QueryAllContainerData(tx), "", "\t")
	commons.Must(err)
	commons.Must(tx.Commit())
	fmt.Println(string(b))
}

func HandleInspectContainer() {
	DB := db.Setup(dbPath)
	tx, err := DB.DB.Begin(true)
	commons.Must(err)
	defer func() {
		_ = tx.Rollback()
		commons.Must(DB.DB.Close())
	}()
	b, err := json.MarshalIndent(DB.QueryContainerData(tx, configuration.ID), "", "\t")
	commons.Must(err)
	commons.Must(tx.Commit())
	fmt.Println(string(b))
}

func HandleCleanupResources() {
	DB := db.Setup(dbPath)
	tx, err := DB.DB.Begin(true)
	commons.Must(err)
	defer func() {
		_ = tx.Rollback()
		commons.Must(DB.DB.Close())
	}()
	configuration = DB.QueryContainerData(tx, configuration.ID).Configuration
	DB.DeleteContainerData(tx, configuration.ID)
	CleanupContainerResources()
	commons.Must(tx.Commit())
}

func HandleCreateBridgeNetwork() {
	DB := db.Setup(dbPath)
	tx, err := DB.DB.Begin(true)
	commons.Must(err)
	defer func() {
		_ = tx.Rollback()
		commons.Must(DB.DB.Close())
	}()
	networkData := DB.QueryNetworkData(tx, networkOpts.Name)
	if networkData != nil {
		fmt.Printf("network %s existed", networkOpts.Name)
		return
	}
	CreateBridgeNetwork(networkOpts.Name, networkOpts.Address)
	DB.SaveNetworkData(tx, networkOpts.Name, networkOpts.Address, make(map[string]bool))
	commons.Must(tx.Commit())
}

func HandleDeleteBridgeNetwork() {
	DB := db.Setup(dbPath)
	tx, err := DB.DB.Begin(true)
	commons.Must(err)
	defer func() {
		_ = tx.Rollback()
		commons.Must(DB.DB.Close())
	}()
	networkData := DB.QueryNetworkData(tx, networkOpts.Name)
	if networkData != nil && len(networkData.AttachedContainers) > 0 {
		fmt.Printf("network %s have some container attached to", networkOpts.Name)
		return
	}
	DeleteBridgeNetwork(networkOpts.Name)
	DB.DeleteNetworkData(tx, networkOpts.Name)
	commons.Must(tx.Commit())
}

func HandleListBridgeNetwork() {
	DB := db.Setup(dbPath)
	tx, err := DB.DB.Begin(true)
	commons.Must(err)
	defer func() {
		_ = tx.Rollback()
		commons.Must(DB.DB.Close())
	}()
	b, err := json.MarshalIndent(DB.QueryAllNetworkData(tx), "", "\t")
	commons.Must(err)
	commons.Must(tx.Commit())
	fmt.Println(string(b))
}

func HandleRun() {
	cmd := exec.Cmd{
		Path:   "/proc/self/exe",
		Args:   append([]string{"run", configuration.ID}, os.Args[1:]...),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		SysProcAttr: &syscall.SysProcAttr{
			Pdeathsig:  syscall.SIGTERM,
			Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
			UidMappings: []syscall.SysProcIDMap{
				{
					ContainerID: 0,
					HostID:      configuration.UID,
					Size:        1,
				},
			},
			GidMappings: []syscall.SysProcIDMap{
				{
					ContainerID: 0,
					HostID:      configuration.GID,
					Size:        1,
				},
			},
			Credential: &syscall.Credential{
				Uid: 0,
				Gid: 0,
			},
		},
	}
	commons.Must(cmd.Start())
	func() {
		DB := db.Setup(dbPath)
		tx, err := DB.DB.Begin(true)
		commons.Must(err)
		defer func() {
			_ = tx.Rollback()
			commons.Must(DB.DB.Close())
		}()
		InitCGroup(configuration.Resources.Cgroup, cmd.Process.Pid)
		CGroupLimitCPU(configuration.Resources.Cgroup, configuration.Resources.CPU)
		CGroupLimitMemory(configuration.Resources.Cgroup, configuration.Resources.Memory)
		DB.SaveContainerData(tx, configuration.ID, configuration, cmd.Process.Pid, make(map[string]db.Veth))
		for networkName, network := range configuration.Networks {
			networkData := DB.QueryNetworkData(tx, networkOpts.Name)
			if networkData != nil {
				fmt.Printf("network %s existed", networkOpts.Name)
			} else {
				CreateBridgeNetwork(networkName, network.Address)
				DB.SaveNetworkData(tx, networkName, network.Address, make(map[string]bool))
			}
			veth := CreateVethPair(networkName, cmd.Process.Pid)
			DB.AddContainerToNetwork(tx, networkName, network.Address, configuration.ID, veth, network.CIDR)
		}
		commons.Must(tx.Commit())
	}()
	_ = cmd.Wait()
	// cleanup container
	func() {
		DB := db.Setup(dbPath)
		tx, err := DB.DB.Begin(true)
		commons.Must(err)
		defer func() {
			_ = tx.Rollback()
			commons.Must(DB.DB.Close())
		}()
		containerData := DB.QueryContainerData(tx, configuration.ID)
		for network, _ := range containerData.Configuration.Networks {
			DB.RemoveContainerFromNetwork(tx, network, configuration.ID)
		}
		DB.DeleteContainerData(tx, configuration.ID)
		commons.Must(tx.Commit())
	}()
}

func HandleCLI() {
	rand.Seed(time.Now().UnixNano())
	ParseArgs()
	switch choiceOpts.Opt {
	case "list":
		HandleListContainer()
		goto FINISH
	case "inspect":
		HandleInspectContainer()
		goto FINISH
	case "cleanup":
		commons.MustBeExecutedByRoot()
		HandleCleanupResources()
		goto FINISH
	case "network":
		switch networkOpts.Action {
		case "create":
			HandleCreateBridgeNetwork()
			goto FINISH
		case "delete":
			HandleDeleteBridgeNetwork()
			goto FINISH
		case "list":
			HandleListBridgeNetwork()
			goto FINISH
		}
		goto FINISH
	case "run":
		commons.MustBeExecutedByRoot()
		HandleRun()
		goto FINISH
	}
FINISH:
}
