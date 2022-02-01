package main

import (
	"experiment_lwc/commons"
	"fmt"
	"github.com/vishvananda/netlink"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type Mount struct {
	Source string
	Target string
	Fs     string
	Flags  int
	Data   string
}

func mount(mounts []Mount) {
	for _, mnt := range mounts {
		target := filepath.Join(*rootfs, mnt.Target)
		commons.Must(syscall.Mount(mnt.Source, target, mnt.Fs, uintptr(mnt.Flags), mnt.Data))
	}
}

func pivotAndMountRootFs(rootFs string) {
	pivotDir := "pivot_fs"
	pivotPath := path.Join(rootFs, pivotDir)
	commons.Must(syscall.Mount(rootFs, rootFs, "bind", syscall.MS_BIND|syscall.MS_REC, ""))
	commons.Must(os.MkdirAll(pivotPath, 0700))
	commons.Must(syscall.PivotRoot(rootFs, pivotPath))
	commons.Must(os.Chdir("/"))
	commons.Must(syscall.Unmount(pivotDir, syscall.MNT_DETACH))
	commons.Must(os.RemoveAll(pivotDir))
}

func createBridgeNetwork(bridgeName string, bridgeAddress string) {
	_, err := net.InterfaceByName(bridgeName)
	if err == nil {
		return
	}
	if !strings.Contains(err.Error(), "no such network interface") {
		panic(err)
	}
	linkAttrs := &netlink.LinkAttrs{
		Name:   bridgeName,
		TxQLen: -1,
	}
	bridge := &netlink.Bridge{
		LinkAttrs: *linkAttrs,
	}
	commons.Must(netlink.LinkAdd(bridge))
	address, err := netlink.ParseAddr(bridgeAddress)
	commons.Must(err)
	commons.Must(netlink.AddrAdd(bridge, address))
	commons.Must(netlink.LinkSetUp(bridge))
}

func createVethPair(bridgeName string, pid int) {
	bridge, err := netlink.LinkByName(bridgeName)
	commons.Must(err)
	parentName := fmt.Sprintf("veth_%s", commons.StringRandom(8, commons.Lowercase+commons.Numeric))
	peerName := fmt.Sprintf("veth_%s", commons.StringRandom(8, commons.Lowercase+commons.Numeric))
	linkAttrs := &netlink.LinkAttrs{
		Name:        parentName,
		TxQLen:      -1,
		MasterIndex: bridge.Attrs().Index,
	}
	vethPair := &netlink.Veth{
		LinkAttrs: *linkAttrs,
		PeerName:  peerName,
	}
	err = netlink.LinkAdd(vethPair)
	if err != nil && strings.Contains(err.Error(), "file exists") {
		commons.Must(netlink.LinkDel(vethPair))
		commons.Must(netlink.LinkAdd(vethPair))
	}
	peer, err := netlink.LinkByName(peerName)
	commons.Must(err)
	commons.Must(netlink.LinkSetNsPid(peer, pid))
	commons.Must(netlink.LinkSetUp(vethPair))
}

func SetupBridge(bridgeName string, bridgeAddress string, pid int) {
	createBridgeNetwork(bridgeName, bridgeAddress)
	createVethPair(bridgeName, pid)
}

func setupNetworkInterface() {
	start := time.Now()
	var err error
	var link netlink.Link
	done := false
	for {
		if time.Since(start) > 5*time.Second {
			err = fmt.Errorf("failed to find veth interface")
		}
		commons.Must(err)
		linkList, err := netlink.LinkList()
		commons.Must(err)
		for _, temp := range linkList {
			if temp.Type() == "veth" {
				link = temp
				done = true
			}
		}
		if done {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	lo, err := netlink.LinkByName("lo")
	commons.Must(err)
	commons.Must(netlink.LinkSetUp(lo))
	addr, err := netlink.ParseAddr(*containerIp)
	commons.Must(err)
	commons.Must(netlink.AddrAdd(link, addr))
}

func SpawnContainer() {
	mount([]Mount{
		{
			Source: "proc",
			Target: "/proc",
			Fs:     "proc",
			Flags:  syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV,
		},
	})
	pivotAndMountRootFs(*rootfs)
	commons.Must(os.Chdir(*chdir))
	commons.Must(syscall.Sethostname([]byte(*hostname)))
	setupNetworkInterface()
	commons.Must(syscall.Exec("/bin/sh", []string{"sh"}, os.Environ()))
}
