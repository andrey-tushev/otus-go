package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	vars := make(Environment)
	for _, f := range files {
		// Если файл нулевого размера
		ev := EnvValue{
			NeedRemove: true,
		}

		// Если файл не нулевого размера
		if f.Size() > 0 {
			value, err := fetchVar(dir + "/" + f.Name())
			if err != nil {
				return Environment{}, err
			}

			ev = EnvValue{
				Value:      value,
				NeedRemove: false,
			}
		}

		vars[f.Name()] = ev
	}

	return vars, nil
}

func fetchVar(filePath string) (string, error) {
	file, err := os.Open(filePath)
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
