/*
Copyright 2016 The ContainerOps Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/urfave/cli"
	"gopkg.in/macaron.v1"

	"github.com/liangchenye/update-service/utils"
	_ "github.com/liangchenye/update-service/utils/km/local"
	_ "github.com/liangchenye/update-service/utils/protocal/appV1"
	_ "github.com/liangchenye/update-service/utils/storage/local"
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
			Name:  "storage",
			Value: "local://tmp/dockyard-updater-server-storage",
			Usage: "the storage database",
		},
		cli.StringFlag{
			Name:  "keymanager",
			Value: "local://tmp/dockyard-updater-server-keymanager",
			Usage: "the key manager url",
		},
	},
}

func runUpdateServer(c *cli.Context) error {
	m := macaron.New()

	for _, item := range []string{"keymanager", "storage"} {
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
