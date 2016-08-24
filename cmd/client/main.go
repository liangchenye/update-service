package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "uc"
	app.Usage = "Updater Service client"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		addCommand,
		removeCommand,
		listCommand,
		pushCommand,
		pullCommand,
	}

	app.Run(os.Args)
}
