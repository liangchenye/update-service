package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/urfave/cli"
	"gopkg.in/macaron.v1"

	"github.com/liangchenye/update-service/utils"
)

var webCommand = cli.Command{
	Name:        "web",
	Usage:       "Update Server",
	Description: "Update Server stores the signatured meta data.",
	Action:      runUpdateServer,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "address",
			Value: "0.0.0.0",
			Usage: "web service listen ip, default is 0.0.0.0; if listen with Unix Socket, the value is sock file path.",
		},
		cli.StringFlag{
			Name:  "listen-mode",
			Value: "http",
			Usage: "web service listen mode, default is http.",
		},
		cli.IntFlag{
			Name:  "port",
			Value: 1234,
			Usage: "web service listen at port 80; if run with https will be 443.",
		},
		cli.StringFlag{
			Name:  "storage-uri",
			Value: "/tmp/dockyard-updater-server-storage",
			Usage: "the storage database",
		},
		cli.StringFlag{
			Name:  "keymanager-mode",
			Value: "peruser",
			Usage: "the key manager mode",
		},
		cli.StringFlag{
			Name:  "keymanager-uri",
			Value: "/tmp/dockyard-updater-server-keymanager",
			Usage: "the key manager url",
		},
	},
}

func runUpdateServer(c *cli.Context) error {
	m := macaron.New()

	for _, item := range []string{"keymanager-mode", "keymanager-uri", "storage-uri"} {
		utils.SetSetting(item, c.String(item))
	}

	SetRouters(m)

	switch c.String("listen-mode") {
	case "http":
		listenaddr := fmt.Sprintf("%s:%d", c.String("address"), c.Int("port"))
		fmt.Printf("Start listen to :%s\n", listenaddr)
		if err := http.ListenAndServe(listenaddr, m); err != nil {
			fmt.Printf("Start Update Server http mode error: %v\n", err.Error())
			return err
		}
		break
	default:
		break
	}

	return nil
}

func main() {
	app := cli.NewApp()

	app.Name = "upserver"
	app.Usage = "Update Server"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		webCommand,
	}

	app.Run(os.Args)
}
