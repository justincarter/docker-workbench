package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
	"github.com/justincarter/docker-workbench/machine"
)

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
	machine.EvalHint(name, false)
	printWorkbenchInfo("*", name)

	return nil
}

// Up command
func Up(c *cli.Context) error {

	// get name from the current working directory
	workdir, _ := os.Getwd()
	name := filepath.Base(workdir)

	if machine.Exists(name) {
		machine.Start(name)
		machine.EvalHint(name, true)
		printWorkbenchInfo("*", name)
	} else {
		// get name from the parent of the current working directory
		app := name
		name := filepath.Base(filepath.Dir(workdir))

		if machine.Exists(name) {
			machine.Start(name)
			machine.EvalHint(name, true)

			fmt.Println("\nStart the application:")
			fmt.Println("docker-compose up")
			printWorkbenchInfo(app, name)
		} else {
			fmt.Printf("Workbench machine '%s' not found.\n", app)
			os.Exit(1)
		}
	}
	return nil
}

// printWorkbenchInfo prints the application URL using the given app name and workbench machine IP
func printWorkbenchInfo(app, name string) {
	ip, ok := machine.IP(name)
	if ok == true {
		fmt.Println("\nBrowse the workbench using:")
		fmt.Printf("http://%s.%s.nip.io/\n", app, ip)
	} else {
		fmt.Println("\nCould not find the IP address for this workbench")
		os.Exit(1)
	}
}
