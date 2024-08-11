package main

import (
	"os"
	"strings"
)

type Environment map[string]EnvValue

type EnvValue struct {
	Value      string
	NeedRemove bool
}

func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := dir + "/" + file.Name()
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		if len(fileData) == 0 {
			env[file.Name()] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
			continue
		}

		lines := strings.SplitN(string(fileData), "\n", 2)
		value := strings.ReplaceAll(lines[0], "\x00", "\n")
		value = strings.TrimRight(value, " \t\n")

		env[file.Name()] = EnvValue{
			Value:      value,
			NeedRemove: false,
		}
	}

	return env, nil
}
