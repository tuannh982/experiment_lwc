package main

import (
	"os"
)

const dbPath = "database.db"

func main() {
	if os.Args[0] == "run" { // child process, /proc/self/exe
		HandleChildProcess()
	} else {
		HandleCLI()
	}
}
