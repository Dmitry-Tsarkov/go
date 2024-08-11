package main

import (
	"errors"
	"os"
	"os/exec"
)

func RunCmd(cmd []string, env Environment) (string, int) {
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	for key, value := range env {
		if value.NeedRemove {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, value.Value)
		}
	}

	out, err := command.CombinedOutput()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return string(out), exitError.ExitCode()
		}
		return string(out), 1
	}

	return string(out), 0
}
