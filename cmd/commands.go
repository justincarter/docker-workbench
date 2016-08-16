package cmd

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/justincarter/docker-workbench/machine"
	"github.com/urfave/cli"
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
	{
		Name:   "proxy",
		Usage:  "Start a reverse proxy to the app in the current directory",
		Action: Proxy,
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

// Proxy command
func Proxy(c *cli.Context) error {

	// get name from the current working directory
	workdir, _ := os.Getwd()
	name := filepath.Base(workdir)
	app := ""

	if machine.Exists(name) {
		fmt.Printf("\nCould not find the app to proxy for Workbench machine '%s'", name)
		os.Exit(1)
	} else {
		// get name from the parent of the current working directory
		app = name
		name = filepath.Base(filepath.Dir(workdir))

		if machine.Exists(name) {
			ip, ok := machine.IP(name)
			if ok == true {
				proxy := httputil.NewSingleHostReverseProxy(&url.URL{
					Scheme: "http",
					Host:   fmt.Sprintf("%s.%s.nip.io", app, ip),
				})
				fmt.Println("Starting reverse proxy on port 9999...")

				ifaces, err := net.Interfaces()
				if err != nil {
					fmt.Println("\nCould not find local network interfaces")
					os.Exit(1)
				}
				fmt.Printf("Listening on:\n\n")
				for _, i := range ifaces {
					addrs, _ := i.Addrs()
					for _, addr := range addrs {
						var ip string
						switch v := addr.(type) {
						case *net.IPNet:
							ip = v.IP.String()
						case *net.IPAddr:
							ip = v.IP.String()
						}
						// output valid local IPv4 addresses, excluding loopbacks and docker machine default interface
						if machine.ValidIPv4(ip) && ip != "127.0.0.1" && ip != "192.168.99.1" && strings.Split(ip, ".")[0] != "169" {
							fmt.Printf("http://%s.%s.nip.io:9999/\n", app, ip)
						}
					}
				}
				fmt.Println("\nPress Ctrl-C to terminate proxy")

				l, err := net.Listen("tcp4", ":9999")
				if err != nil {
					log.Fatal(err)
				}
				log.Fatal(http.Serve(l, proxy))

			} else {
				fmt.Println("\nCould not find the IP address for this workbench")
				os.Exit(1)
			}

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
