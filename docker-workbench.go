package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/justincarter/docker-workbench/cmd"
)

const version = "0.5"

func main() {

	if err := flightCheck(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cli.AppHelpTemplate = templateAppHelp
	cli.CommandHelpTemplate = templateCommandHelp
	cli.VersionPrinter = cmd.Version

	app := cli.NewApp()
	app.Name = "docker-workbench"
	app.Version = version
	app.Usage = "Provision a Docker Workbench for use with docker-machine and docker-compose"

	app.CommandNotFound = cmd.NotFound
	app.Commands = cmd.Commands

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

func flightCheck() error {

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
	if _, err := exec.LookPath("VBoxManage"); err != nil {
		return fmt.Errorf("docker-workbench: VBoxManage was not found. Make sure you have installed VirtualBox")
	}

	return nil
}
