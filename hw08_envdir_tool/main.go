package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
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

	vars, err := getVars(dir)
	if err != nil {
		log.Fatal(err)
	}

	for name, value := range vars {
		_, exists := os.LookupEnv(name)
		if exists {
			os.Unsetenv(name)
		}
		os.Setenv(name, value)
	}

	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func getVars(dir string) (map[string]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	vars := make(map[string]string)
	for _, f := range files {
		value, err := getVar(dir + "/" + f.Name())
		if err != nil {
			return nil, err
		}

		vars[f.Name()] = value
	}

	return vars, nil
}

func getVar(name string) (string, error) {
	file, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	value := scanner.Text()

	return trim(value), nil
}

func trim(value string) string {
	value = strings.TrimRight(value, " \t")
	value = strings.ReplaceAll(value, "\x00", "\n")
	return value
}
