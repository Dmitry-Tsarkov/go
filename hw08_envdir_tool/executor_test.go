package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestRunCmd(t *testing.T) {
	envVars := Environment{
		"FOO": EnvValue{Value: "123", NeedRemove: false},
		"BAR": EnvValue{Value: "value", NeedRemove: false},
	}

	cmdArgs := []string{"sh", "-c", "echo $FOO $BAR"}

	outputBuffer := &bytes.Buffer{}
	command := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	command.Env = os.Environ()
	for key, value := range envVars {
		if value.NeedRemove {
			command.Env = removeEnv(command.Env, key)
		} else {
			command.Env = append(command.Env, key+"="+value.Value)
		}
	}

	command.Stdout = outputBuffer
	command.Stderr = outputBuffer

	if err := command.Run(); err != nil {
		t.Fatal(err)
	}

	expected := "123 value\n"
	if outputBuffer.String() != expected {
		t.Errorf("Expected %s, got %s", expected, outputBuffer.String())
	}
}

func removeEnv(env []string, key string) []string {
	newEnv := []string{}
	for _, v := range env {
		if !strings.HasPrefix(v, key+"=") {
			newEnv = append(newEnv, v)
		}
	}
	return newEnv
}
