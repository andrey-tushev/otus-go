package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from envs.
func RunCmd(command string, args []string, envs Environment) (returnCode int) {
	for name, env := range envs {
		if env.NeedRemove {
			os.Unsetenv(name)
		} else {
			_, exists := os.LookupEnv(name)
			if exists {
				os.Unsetenv(name)
			}
			os.Setenv(name, env.Value)
		}
	}

	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		exitError, isExitError := err.(*exec.ExitError) // nolint:errorlint
		if isExitError {
			return exitError.ExitCode()
		}

		return 1
	}

	return 0
}
