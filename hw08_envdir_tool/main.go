package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("not enough arguments")
	}

	dir := os.Args[1]
	command := os.Args[2]

	var args []string
	if len(os.Args) > 3 {
		args = os.Args[3:]
	}

	envs, err := ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	retCode := RunCmd(command, args, envs)
	os.Exit(retCode)
}
