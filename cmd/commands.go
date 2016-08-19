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

type workbench struct {
	app  string
	name string
}

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
	machine.EvalHint(name, false)
	printWorkbenchInfo("*", name)

	return nil
}

// Up command
func Up(c *cli.Context) error {
	w, err := getWorkbenchContext(false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	machine.Start(w.name)
	machine.EvalHint(w.name, true)
	if w.app != "*" {
		fmt.Println("\nStart the application:")
		fmt.Println("docker-compose up")
	}
	printWorkbenchInfo(w.app, w.name)

	return nil
}

// Proxy command
func Proxy(c *cli.Context) error {
	w, err := getWorkbenchContext(true)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ip, ok := machine.IP(w.name)
	if ok == true {
		fmt.Printf("Starting reverse proxy on port %s...\n", proxyPort)
		ips, err := getProxyIPs()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Listening on:\n\n")
		for _, thisip := range ips {
			fmt.Printf("http://%s.%s.nip.io:%s/\n", w.app, thisip, proxyPort)
		}
		fmt.Println("\nPress Ctrl-C to terminate proxy")
		startProxy(w, ip, proxyPort)
	} else {
		fmt.Println("\nCould not find the IP address for this workbench. Have you run docker-workbench up?")
		os.Exit(1)
	}

	return nil
}

// getWorkbenchContext finds the application name and workbench machine name  from the current directory
func getWorkbenchContext(requireApp bool) (w workbench, err error) {
	err = nil
	// get name from the current working directory
	workdir, _ := os.Getwd()
	w = workbench{
		app:  "*",
		name: filepath.Base(workdir),
	}

	if !machine.Exists(w.name) {
		// get name from the parent of the current working directory
		w.app = w.name
		w.name = filepath.Base(filepath.Dir(workdir))

		if !machine.Exists(w.name) {
			err = fmt.Errorf("Workbench machine '%s' not found.", w.app)
		}
	}
	if requireApp == true && w.app == "*" {
		err = fmt.Errorf("\nCould not find the app to proxy for Workbench machine '%s'. Try running from an app directory?", w.name)
	}
	return
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

// getProxyIPs returns a slice of IP address strings that should be browsable when using the Proxy command
func getProxyIPs() ([]string, error) {
	var e error

	ifaces, err := net.Interfaces()
	if err != nil {
		e = fmt.Errorf("\nCould not find local network interfaces")
	}

	ips := getIPsFromIfaces(ifaces)
	if len(ips) == 0 {
		e = fmt.Errorf("\nCould not find local network interfaces")
	}

	return ips, e
}

func getIPsFromIfaces(ifaces []net.Interface) []string {
	ips := []string{}
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
			if validProxyIP(ip) {
				ips = append(ips, ip)
			}
		}
	}
	return ips
}

func validProxyIP(ip string) bool {
	// disallow non-IPv4 addresses, loopback interfaces and docker machine default interface
	if !machine.ValidIPv4(ip) || ip == "127.0.0.1" || ip == "192.168.99.1" || strings.Split(ip, ".")[0] == "169" {
		return false
	}
	return true
}

func startProxy(w workbench, ip, port string) {
	l, err := net.Listen("tcp4", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s.%s.nip.io", w.app, ip),
	})
	log.Fatal(http.Serve(l, proxy))
}
