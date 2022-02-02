package main

import (
	"experiment_lwc/commons"
	"experiment_lwc/config"
	"fmt"
	cp "github.com/otiai10/copy"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func copyRootFs(rootfs string, path string) {
	commons.Must(os.MkdirAll(path, 0755))
	commons.Must(cp.Copy(rootfs, path))
}

func SpawnContainer() {
	newRootFs := filepath.Join("containers", configuration.ID)
	copyRootFs(configuration.RootFs, newRootFs)
	mount(newRootFs, append(configuration.Mounts, []config.Mount{
		{
			Source: "proc",
			Target: "/proc",
			Fs:     "proc",
			Flags:  syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV,
		},
	}...))
	pivotAndMountRootFs(newRootFs)
	commons.Must(os.Chdir(configuration.Chdir))
	commons.Must(syscall.Sethostname([]byte(configuration.Hostname)))
}

func CleanupContainerResources() {
	// will not remove bridge networks
	cmd := exec.Command("cgdelete", fmt.Sprintf("cpu,memory:/%s", configuration.Resources.Cgroup))
	commons.Must(cmd.Run())
}
