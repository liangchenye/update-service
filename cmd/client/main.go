package main

import (
	"os"

	"github.com/urfave/cli"

	_ "github.com/liangchenye/update-service/cmd/client/utils/repo/appV1"
)

func main() {
	app := cli.NewApp()

	app.Name = "duc"
	app.Usage = "Dockyard Updater client"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		initCommand,
		addCommand,
		removeCommand,
		listCommand,
		pushCommand,
		pullCommand,
	}

	app.Run(os.Args)
}
