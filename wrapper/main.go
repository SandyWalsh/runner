package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func addProcessToCgroup(fn string, pid int) {
	file, err := os.OpenFile(fn, os.O_WRONLY, 0644)
	if err != nil {
		log.Println("unable to open", fn, "for writing :", err)
		os.Exit(1)
	}
	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("%d", pid)); err != nil {
		log.Println("failed to add pid to cgroup: ", err)
		os.Exit(1)
	}
}

func main() {
	// wrapper procfile cmd args...
	if len(os.Args) < 3 {
		log.Fatalln("too few args to wrapper")
		os.Exit(1)
	}
	pfile := os.Args[1]
	cmd := os.Args[2]
	args := os.Args[3:]

	proc := fmt.Sprintf("%s/%s", pfile, "cgroup.procs")

	addProcessToCgroup(proc, os.Getpid())

	b, err := exec.LookPath(cmd)
	if err != nil {
		log.Fatalln("cannot find path for cmd:", err)
		os.Exit(1)
	}

	err = syscall.Exec(b, args, os.Environ())
	if err != nil {
		log.Fatalln("unable to launch program:", err)
		os.Exit(1)
	}
}
