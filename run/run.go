package run

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Run is helper for running a command with known arguments from a slice of strings
func Run(command string, args []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ReadOutput is helper for running a command and returning its output
func ReadOutput(command, argstr string, vars ...interface{}) ([]byte, error) {
	args := fmt.Sprintf(argstr, vars...)
	cmd := exec.Command(command, strings.Fields(args)...)
	return cmd.Output()
}
