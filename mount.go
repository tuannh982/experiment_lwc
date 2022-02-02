package main

import (
	"experiment_lwc/commons"
	"experiment_lwc/config"
	"os"
	"path"
	"path/filepath"
	"syscall"
)

func mount(rootfs string, mounts []config.Mount) {
	for _, mnt := range mounts {
		target := filepath.Join(rootfs, mnt.Target)
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
