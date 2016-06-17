package main

import (
	"fmt"
	"os"

	"github.com/justincarter/docker-workbench/cmd"
	"github.com/urfave/cli"
)

const version = "0.5"

func main() {

	if err := cmd.FlightCheck(); err != nil {
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
