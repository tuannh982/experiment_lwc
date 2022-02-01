package main

import (
	"experiment_lwc/commons"
	"flag"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

var (
	rootfs        = flag.String("rootfs", "./rootfs", "rootfs")
	chdir         = flag.String("chdir", "/", "chdir")
	uid           = flag.Int("uid", 1000, "UID")
	gid           = flag.Int("gid", 1000, "GID")
	hostname      = flag.String("hostname", "container", "hostname")
	cgroup        = flag.String("cgroup", "isolator_container", "cgroup name")
	cpu           = flag.Float64("cpu", -1, "cpu")
	memory        = flag.Int("memory", -1, "memory")
	bridgeNetwork = flag.String("bridge", "bridge", "bridge network name")
	bridgeAddress = flag.String("bridge-addr", "10.10.10.1/24", "bridge network")
	containerIp   = flag.String("ip", "10.10.10.2/24", "container CIDR address")
)

func getExecUID() int {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()
	commons.Must(err)
	i, err := strconv.Atoi(string(output[:len(output)-1]))
	commons.Must(err)
	return i
}

func main() {
	if getExecUID() != 0 {
		panic("this program must be executed as root")
	}
	flag.Parse()
	if os.Args[0] == "spawn-container" {
		SpawnContainer()
		os.Exit(0)
	}
	cmd := exec.Cmd{
		Path:   "/proc/self/exe",
		Args:   append([]string{"spawn-container"}, os.Args[1:]...),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		SysProcAttr: &syscall.SysProcAttr{
			Pdeathsig:  syscall.SIGTERM,
			Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
			UidMappings: []syscall.SysProcIDMap{
				{
					ContainerID: 0,
					HostID:      *uid,
					Size:        1,
				},
			},
			GidMappings: []syscall.SysProcIDMap{
				{
					ContainerID: 0,
					HostID:      *gid,
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
	InitCGroup(*cgroup, cmd.Process.Pid)
	CGroupLimitCPU(*cgroup, *cpu)
	CGroupLimitMemory(*cgroup, *memory)
	SetupBridge(*bridgeNetwork, *bridgeAddress, cmd.Process.Pid)
	commons.Must(cmd.Wait())
}
