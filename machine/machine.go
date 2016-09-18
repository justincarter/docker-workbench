package machine

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/justincarter/docker-workbench/run"
)

// Machine represents a docker machine
type Machine struct {
	Name string
}

// Create the docker machine
func (m *Machine) Create() {

	// default configuration using docker-machine environment variables
	env := map[string]string{
		"VIRTUALBOX_CPU_COUNT":   "2",
		"VIRTUALBOX_DISK_SIZE":   "60000",
		"VIRTUALBOX_MEMORY_SIZE": "2048",
		"VIRTUALBOX_NO_SHARE":    "true",
	}
	for k, v := range env {
		// set defaults only if the variable is currently unset
		if os.Getenv(k) == "" {
			os.Setenv(k, v)
		}
	}

	// create the machine
	args := fmt.Sprintf("create --driver virtualbox %s", m.Name)
	cmd := exec.Command("docker-machine", strings.Fields(args)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("docker-workbench: docker-machine create failed.")
		os.Exit(1)
	}

}

// EvalEnv sets docker environment variables
func (m *Machine) EvalEnv() {
	out, _ := run.Output("docker-machine", "env", m.Name, "--shell=bash")
	env := parseEnvOutput(out)
	for k, v := range env {
		os.Setenv(k, v)
	}
}

// PrintEvalHint shows a hint about running docker env if required
func (m *Machine) PrintEvalHint(checkenv bool) {
	showhint := true
	if checkenv == true && os.Getenv("DOCKER_MACHINE_NAME") == m.Name {
		showhint = false
	}
	if showhint == true {
		fmt.Println("\nRun the following command to set this machine as your default:")
		fmt.Printf("eval \"$(docker-machine env %s)\"\n", m.Name)
	}
}

// Exists checks if a VM exists
func (m *Machine) Exists() bool {
	out, _ := run.Output(VBoxManagePath(), "list", "vms")
	re := regexp.MustCompile("(?mi)^\"" + m.Name + "\"")
	return re.Match(out)
}

// IP returns the IP address of the docker machine
func (m *Machine) IP() (ip string, success bool) {
	out, _ := run.Output("docker-machine", "ip", m.Name)
	ip = strings.Split(string(out), "\n")[0]
	success = ValidIPv4(ip)
	return
}

// ShareFolder adds a /workbench shared folder to the VM
func (m *Machine) ShareFolder(folder string) {
	args := []string{"sharedfolder", "add", m.Name, "--name", "workbench", "--hostpath", folder}
	run.Run(VBoxManagePath(), args...)
}

// SSH into the docker machine to run a command
func (m *Machine) SSH(command string) {
	run.Run("docker-machine", "ssh", m.Name, command)
}

// Start the docker machine
func (m *Machine) Start() {
	run.Run("docker-machine", "start", m.Name)
}

// Stop the docker machine
func (m *Machine) Stop() {
	run.Run("docker-machine", "stop", m.Name)
}

// ValidIPv4 returns true for valid IPv4 addresses
func ValidIPv4(ip string) bool {
	// validate IP address
	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	return re.Match([]byte(ip))
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

// parseEnvOutput parses the output from `docker-machine env` and returns a map
func parseEnvOutput(output []byte) map[string]string {
	env := make(map[string]string)
	for _, line := range strings.Split(string(output), "\n") {
		re, _ := regexp.Compile(`export (.*?)="(.*)"`)
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			env[matches[1]] = matches[2]
		}
	}
	return env
}
