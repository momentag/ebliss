package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
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
	}
	app.Version = "v0.1"
	app.EnableBashCompletion = true

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
