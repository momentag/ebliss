package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/momentag/ebliss/auth"
	"github.com/urfave/cli/v2"
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
		Commands: []*cli.Command{
			{
				Name:        "auth",
				Usage:       "provides authentication utilities, including a server",
				Subcommands: auth.Commands(),
			},
		},
	}
	app.Version = version
	app.EnableBashCompletion = true
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
