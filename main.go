package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"

	_ "github.com/hashicorp/hcl/v2"

	_ "github.com/momentag/ebliss/config"
	"github.com/momentag/ebliss/gin"
)

var (
	version = fmt.Sprintf("v0.1-%v", uuid.New().String()[0:8])
)

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "prints the current installed version",
	}
	app := &cli.App{
		Name:     "ebliss",
		HelpName: "ebliss",
		Usage:    "An event-driven enterprise resource planning application",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "friend-addr",
				Usage: "address to join a cluster",
			},
		},
		Commands: []*cli.Command{
			gin.Command(),
		},
	}
	app.Version = version
	app.EnableBashCompletion = true
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
