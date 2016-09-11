package run

import (
	"os"
	"os/exec"
)

// Run is helper for running a command with a variable number of string arguments
func Run(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Output is helper for running a command with a variable number of string arguments and returning its output
func Output(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	return cmd.Output()
}
