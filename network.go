package main

import (
	"experiment_lwc/commons"
	"fmt"
	"github.com/vishvananda/netlink"
	"net"
	"strings"
	"time"
)

func CreateBridgeNetwork(bridgeName string, bridgeAddress string) {
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

func DeleteBridgeNetwork(bridgeName string) {
	_, err := net.InterfaceByName(bridgeName)
	if err != nil && !strings.Contains(err.Error(), "no such network interface") {
		panic(err)
	}
	linkAttrs := &netlink.LinkAttrs{
		Name:   bridgeName,
		TxQLen: -1,
	}
	bridge := &netlink.Bridge{
		LinkAttrs: *linkAttrs,
	}
	commons.Must(netlink.LinkDel(bridge))
}

func CreateVethPair(bridgeName string, pid int) string {
	bridge, err := netlink.LinkByName(bridgeName)
	commons.Must(err)
	hostVethName := fmt.Sprintf("veth_%s", commons.StringRandom(8, commons.Lowercase+commons.Numeric))
	containerVethName := fmt.Sprintf("ceth_%s", commons.StringRandom(8, commons.Lowercase+commons.Numeric))
	linkAttrs := &netlink.LinkAttrs{
		Name:   hostVethName,
		TxQLen: -1,
	}
	vethPair := &netlink.Veth{
		LinkAttrs: *linkAttrs,
		PeerName:  containerVethName,
	}
	commons.Must(netlink.LinkAdd(vethPair))
	hostLink, err := netlink.LinkByName(hostVethName)
	commons.Must(err)
	peerLink, err := netlink.LinkByName(containerVethName)
	commons.Must(err)
	commons.Must(netlink.LinkSetUp(hostLink))
	commons.Must(netlink.LinkSetNsPid(peerLink, pid))
	commons.Must(netlink.LinkSetMaster(hostLink, bridge))
	commons.Must(err)
	return containerVethName
}

func SetupContainerNetworkInterface(vethCidrMap map[string]string) {
	nLinks := len(vethCidrMap)
	start := time.Now()
	var err error
	var links []netlink.Link
	for {
		temp := make([]netlink.Link, 0)
		if time.Since(start) > 10*time.Second {
			err = fmt.Errorf("failed to find veth interface")
		}
		commons.Must(err)
		linkList, err := netlink.LinkList()
		commons.Must(err)
		for _, link := range linkList {
			if link.Type() == "veth" {
				temp = append(temp, link)
			}
		}
		if len(temp) == nLinks {
			links = temp
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	lo, err := netlink.LinkByName("lo")
	commons.Must(err)
	commons.Must(netlink.LinkSetUp(lo))
	for _, link := range links {
		cidr := vethCidrMap[link.Attrs().Name]
		addr, err := netlink.ParseAddr(cidr)
		commons.Must(err)
		commons.Must(netlink.AddrAdd(link, addr))
		commons.Must(netlink.LinkSetUp(link))
	}
}
