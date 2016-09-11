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

// Workbench holds information about a workbench and its app
type Workbench struct {
	App  string
	Name string
}

// NewWorkbench creates a new workbench object
func NewWorkbench(requireApp bool) (*Workbench, error) {
	var err error
	// get name from the current working directory
	workdir, _ := os.Getwd()
	w := &Workbench{
		App:  "*",
		Name: filepath.Base(workdir),
	}

	if !machine.Exists(w.Name) {
		// get name from the parent of the current working directory
		w.App = w.Name
		w.Name = filepath.Base(filepath.Dir(workdir))

		if !machine.Exists(w.Name) {
			err = fmt.Errorf("Workbench machine '%s' not found.", w.App)
		}
	}
	if requireApp == true && w.App == "*" {
		err = fmt.Errorf("\nCould not find the app to proxy for Workbench machine '%s'. Try running from an app directory?", w.Name)
	}
	return w, err
}

// PrintWorkbenchInfo prints the application URL using the app name and machine IP of the workbench
func (w *Workbench) PrintWorkbenchInfo() {
	ip, ok := machine.IP(w.Name)
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
