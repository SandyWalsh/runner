package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func addProcessToCgroup(fn string, pid int) {
	file, err := os.OpenFile(fn, os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("unable to open", fn, "for writing :", err)
	}
	defer file.Close()

	if _, err := fmt.Fprintf(file, "%d", pid); err != nil {
		log.Fatalln("failed to add pid to cgroup: ", err)
	}
}

func main() {
	// args[0] = program
	// args[1] = cgroup.procs filename
	// args[2] = cmd
	// args[3:] = args
	if len(os.Args) < 3 {
		log.Fatalln("too few args to wrapper")
	}
	pfile := os.Args[1]
	cmd := os.Args[2]
	args := os.Args[3:]

	proc := filepath.Join(pfile, "cgroup.procs")

	addProcessToCgroup(proc, os.Getpid())

	b, err := exec.LookPath(cmd)
	if err != nil {
		log.Fatalln("cannot find path for cmd:", err)
	}

	err = syscall.Exec(b, args, os.Environ())
	if err != nil {
		log.Fatalln("unable to launch program:", err)
	}
}
