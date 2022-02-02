package main

import (
	"experiment_lwc/commons"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

func write(path string, content string) {
	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	commons.Must(err)
	_, err = file.WriteString(content)
	commons.Must(err)
	commons.Must(file.Close())
}

func InitCGroup(cgroup string, pid int) {
	resources := []string{"cpu", "memory"}
	for _, resource := range resources {
		dir := fmt.Sprintf("/sys/fs/cgroup/%s/%s/", resource, cgroup)
		commons.Must(os.MkdirAll(dir, 0644))
		write(filepath.Join(dir, "cgroup.procs"), fmt.Sprintf("%d", pid))
	}
}

func CGroupLimitCPU(cgroup string, limit float64) {
	resource := "cpu"
	dir := fmt.Sprintf("/sys/fs/cgroup/%s/%s/", resource, cgroup)
	commons.Must(os.MkdirAll(dir, 0644))
	base := 100000
	var quota int
	if limit < 0 {
		quota = -1
	} else {
		quota = int(limit * float64(base))
	}
	write(filepath.Join(dir, "cpu.cfs_period_us"), fmt.Sprintf("%d", base))
	write(filepath.Join(dir, "cpu.cfs_quota_us"), fmt.Sprintf("%d", quota))
}

func CGroupLimitMemory(cgroup string, limitInBytes int) {
	resource := "memory"
	dir := fmt.Sprintf("/sys/fs/cgroup/%s/%s/", resource, cgroup)
	commons.Must(os.MkdirAll(dir, 0644))
	if limitInBytes < 0 {
		if strconv.IntSize == 64 {
			limitInBytes = 9223372036854771712
		} else if strconv.IntSize == 32 {
			limitInBytes = math.MaxInt - 10240
		}
	}
	write(filepath.Join(dir, "memory.limit_in_bytes"), fmt.Sprintf("%d", limitInBytes))
}
