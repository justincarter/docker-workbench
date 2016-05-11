package machine

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/justincarter/docker-workbench/run"
)

// Create the docker machine
func Create(name string) {

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
	args := fmt.Sprintf("create --driver virtualbox %s", name)
	cmd := exec.Command("docker-machine", strings.Fields(args)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("docker-workbench: docker-machine create failed.")
		os.Exit(1)
	}

}

// EvalEnv sets docker environment variables
func EvalEnv(name string) {
	out, _ := run.ReadOutput("docker-machine", "env %s --shell=bash", name)
	env := parseEnvOutput(out)
	for k, v := range env {
		os.Setenv(k, v)
	}
}

// EvalHint shows a hint about running docker env if required
func EvalHint(name string, checkenv bool) {
	showhint := true
	if checkenv == true && os.Getenv("DOCKER_MACHINE_NAME") == name {
		showhint = false
	}
	if showhint == true {
		fmt.Println("\nRun the following command to set this machine as your default:")
		fmt.Printf("eval \"$(docker-machine env %s)\"\n", name)
	}
}

// IP returns the IP address of the docker machine
func IP(name string) (ip string, success bool) {
	out, _ := run.ReadOutput("docker-machine", "ip %s", name)
	ip = strings.Split(string(out), "\n")[0]
	// validate IP address
	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	success = re.Match([]byte(ip))
	return
}

// SSH into the docker machine to run a command
func SSH(name, command string) {
	run.Run("docker-machine", []string{"ssh", name, command})
}

// Start the docker machine
func Start(name string) {
	run.Run("docker-machine", []string{"start", name})
}

// Stop the docker machine
func Stop(name string) {
	run.Run("docker-machine", []string{"stop", name})
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

// VBoxManage Helpers

// VBoxManagePath returns the path to the VBoxManage executable
func VBoxManagePath() string {
	path := os.Getenv("VBOX_INSTALL_PATH")
	if path == "" {
		path = os.Getenv("VBOX_MSI_INSTALL_PATH")
	}
	if path != "" {
		path += string(os.PathSeparator)
	}
	return path + "VBoxManage"
}

// Exists checks if a VM exists
func Exists(name string) bool {
	out, _ := run.ReadOutput(VBoxManagePath(), "list vms")
	re := regexp.MustCompile("(?mi)^\"" + name + "\"")
	return re.Match(out)
}

// ShareFolder adds a /workbench shared folder to the VM
func ShareFolder(name, folder string) {
	args := []string{"sharedfolder", "add", name, "--name", "workbench", "--hostpath", folder}
	run.Run(VBoxManagePath(), args)
}
