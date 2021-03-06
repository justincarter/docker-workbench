package workbench

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/justincarter/docker-workbench/machine"
)

// Workbench represents a workbench and its app
type Workbench struct {
	machine.Machine
	App string
}

// NewWorkbench creates a new workbench
func NewWorkbench() (*Workbench, error) {
	var err error
	// get name from the current working directory
	workdir, _ := os.Getwd()
	name := filepath.Base(workdir)

	// set up workbench
	w := new(Workbench)
	w.App = "*"
	w.Name = name

	if !w.Exists() {
		// get name from the parent of the current working directory
		name := filepath.Base(filepath.Dir(workdir))

		// set up workbench
		w.App = w.Name
		w.Name = name

		if !w.Exists() {
			err = fmt.Errorf("Workbench machine '%s' not found.", w.App)
		}
	}

	return w, err
}

// PrintWorkbenchInfo prints the application URL using the app name and machine IP of the workbench
func (w *Workbench) PrintWorkbenchInfo() {
	ip, ok := w.IP()
	if ok == true {
		fmt.Println("\nBrowse the workbench using:")
		fmt.Printf("http://%s.%s.nip.io/\n", w.App, ip)
	} else {
		fmt.Println("\nCould not find the IP address for this workbench")
		os.Exit(1)
	}
}

// StartProxy will start a reverse proxy on the given IP address and port number for the workbench
func (w *Workbench) StartProxy(ip, port string) {
	l, err := net.Listen("tcp4", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s.%s.nip.io", w.App, ip),
	})
	log.Fatal(http.Serve(l, proxy))
}

// GetProxyIPs returns a slice of IP address strings that should be browsable when using the Proxy command
func (w *Workbench) GetProxyIPs() ([]string, error) {
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
			appendIPsFromAddr(&ips, addr)
		}
	}
	return ips
}

func appendIPsFromAddr(ips *[]string, addr net.Addr) {
	var ip string
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP.String()
	case *net.IPAddr:
		ip = v.IP.String()
	}
	if validProxyIP(ip) {
		*ips = append(*ips, ip)
	}
}

func validProxyIP(ip string) bool {
	// disallow non-IPv4 addresses, loopback interfaces and docker machine default interface
	if !machine.ValidIPv4(ip) || ip == "127.0.0.1" || ip == "192.168.99.1" || strings.Split(ip, ".")[0] == "169" {
		return false
	}
	return true
}
