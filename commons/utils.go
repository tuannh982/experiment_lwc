package commons

import (
	"os/exec"
	"strconv"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func GetExecUID() int {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()
	Must(err)
	i, err := strconv.Atoi(string(output[:len(output)-1]))
	Must(err)
	return i
}

func MustBeExecutedByRoot() {
	if GetExecUID() != 0 {
		panic("this program must be executed as root")
	}
}
