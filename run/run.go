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

// VBoxManagePath returns the path to the VBoxManage executable
func VBoxManagePath() string {
	path := os.Getenv("VBOX_INSTALL_PATH")
	if path == "" {
		path = os.Getenv("VBOX_MSI_INSTALL_PATH")
	}
	if path != "" && path[len(path)-1:] != string(os.PathSeparator) {
		path += string(os.PathSeparator)
	}
	return path + "VBoxManage"
}
