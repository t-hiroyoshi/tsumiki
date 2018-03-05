package main

import (
	"log"
	"os"

	"github.com/t-hiroyoshi/tsumiki/command"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "Tsumiki"
	app.Usage = "Tsumiki is a container based package manager"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:   "install",
			Usage:  "Install new packages",
			Action: command.InstallAction,
		},
		{
			Name:   "uninstall",
			Usage:  "Uninstall packages",
			Action: command.UninstallAction,
		},
		{
			Name:   "list",
			Usage:  "List all added packages",
			Action: command.ListActions,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
