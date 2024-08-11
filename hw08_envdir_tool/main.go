package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Wrong arguments: %s\n", os.Args[0])
		os.Exit(1)
	}
	envDir := os.Args[1]
	cmdArgs := os.Args[2:]
	envVars, err := ReadDir(envDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading env dir: %v\n", err)
		os.Exit(1)
	}
	output, exitCode := RunCmd(cmdArgs, envVars)
	fmt.Println(output)
	os.Exit(exitCode)
}
