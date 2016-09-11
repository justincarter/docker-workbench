package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/justincarter/docker-workbench/machine"
	"github.com/justincarter/docker-workbench/workbench"
	"github.com/urfave/cli"
)

var proxyPort string

// Commands config
var Commands = []cli.Command{
	{
		Name:   "create",
		Usage:  "Create a new workbench machine in the current directory",
		Action: Create,
	},
	{
		Name:   "up",
		Usage:  "Start the workbench machine and show details",
		Action: Up,
	},
	{
		Name:   "proxy",
		Usage:  "Start a reverse proxy to the app in the current directory",
		Action: Proxy,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "port, p",
				Value:       "8080",
				Usage:       "Port number to start the proxy on",
				Destination: &proxyPort,
			},
		},
	},
}

// FlightCheck helper checks for prerequisite commands
func FlightCheck() error {

	toolbox := []string{"docker", "docker-machine", "docker-compose"}
	missing := []string{}
	for _, c := range toolbox {
		if _, err := exec.LookPath(c); err != nil {
			missing = append(missing, c)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("docker-workbench: %s was not found. Make sure you have installed Docker Toolbox", strings.Join(missing, ", "))
	}
	if _, err := exec.LookPath(machine.VBoxManagePath()); err != nil {
		return fmt.Errorf("docker-workbench: VBoxManage was not found. Make sure you have installed VirtualBox")
	}

	return nil
}

// NotFound command
func NotFound(c *cli.Context, command string) {
	fmt.Printf("docker-workbench: '%s' is not a docker-workbench command. See 'docker-workbench help'.", command)
	os.Exit(1)
}

// Version command
func Version(c *cli.Context) {
	fmt.Printf("v%s", c.App.Version)
}

// Create command
func Create(c *cli.Context) error {

	// get name from the current working directory
	workdir, _ := os.Getwd()
	name := filepath.Base(workdir)

	if !machine.Exists(name) {
		machine.Create(name)
		machine.EvalEnv(name)

		fmt.Println("Configuring bootsync.sh...")
		machine.SSH(name, "sudo echo 'sudo mkdir -p /workbench && sudo mount -t vboxsf -o uid=1000,gid=50 workbench /workbench' >  /tmp/bootsync.sh")
		machine.SSH(name, "sudo cp /tmp/bootsync.sh /var/lib/boot2docker/bootsync.sh")
		machine.SSH(name, "sudo chmod +x /var/lib/boot2docker/bootsync.sh")

		fmt.Println("Installing workbench apps...")
		machine.SSH(name, "docker run -d --restart=always --name=workbench_proxy -p 80:80 -v '/var/run/docker.sock:/tmp/docker.sock:ro' daemonite/workbench-proxy")
		machine.Stop(name)

		fmt.Println("Adding /workbench shared folder...")
		machine.ShareFolder(name, workdir)
	}

	machine.Start(name)
	machine.EvalEnv(name)
	machine.PrintEvalHint(name, false)

	w, _ := workbench.NewWorkbench(false)
	w.PrintWorkbenchInfo()

	return nil
}

// Up command
func Up(c *cli.Context) error {
	w, err := workbench.NewWorkbench(false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	machine.Start(w.Name)
	machine.PrintEvalHint(w.Name, true)
	if w.App != "*" {
		fmt.Println("\nStart the application:")
		fmt.Println("docker-compose up")
	}
	w.PrintWorkbenchInfo()

	return nil
}

// Proxy command
func Proxy(c *cli.Context) error {
	w, err := workbench.NewWorkbench(true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ip, ok := machine.IP(w.Name)
	if !ok {
		fmt.Println("\nCould not find the IP address for this workbench. Have you run docker-workbench up?")
		os.Exit(1)
	}

	fmt.Printf("Starting reverse proxy on port %s...\n", proxyPort)
	ips, err := w.GetProxyIPs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Listening on:\n\n")
	for _, thisip := range ips {
		fmt.Printf("http://%s.%s.nip.io:%s/\n", w.App, thisip, proxyPort)
	}
	fmt.Println("\nPress Ctrl-C to terminate proxy")
	w.StartProxy(ip, proxyPort)

	return nil
}
