package gin

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	eblissConf "github.com/momentag/ebliss/config"
)

func ginServer(service *eblissConf.ServiceConfig) error {
	var router = gin.New()
	addr := strings.ReplaceAll(service.ListenAddr, "http://", "")
	addr = strings.ReplaceAll(addr, "https://", "")
	if err := router.Run(addr); err != nil {
		panic(err)
	}
	return nil
}

func serveService(context *cli.Context) error {
	configFile := context.String("config-file")
	serviceName := context.Args().Get(0)
	if serviceName == "" {
		log.Error().Str("serviceName", serviceName).Msg("service name is missing")
		return nil
	}
	hcl := hclparse.NewParser()
	var config eblissConf.Config
	hclFile, diag := hcl.ParseHCLFile(configFile)
	if diag != nil && diag.HasErrors() {
		panic(diag)
	}
	if hclFile != nil {
		diag = gohcl.DecodeBody(hclFile.Body, nil, &config)
		if diag != nil && diag.HasErrors() {
			panic(diag)
		}
		log.Debug().Str("service", serviceName).Interface("service", config).Msg("starting gin server")
		var service *eblissConf.ServiceConfig
		for _, svc := range config.Services {
			if svc.Name == serviceName {
				service = svc
			}
		}
		if service != nil {
			return ginServer(service)
		}
	}
	return nil
}

func commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "serve",
			Aliases: []string{"s"},
			Action: func(context *cli.Context) error {
				return serveService(context)
			},
		},
	}
}

func Command() *cli.Command {
	return &cli.Command{
		Name:  "gin",
		Usage: "provides gin servers for various services",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config-file",
				Usage:    "configuration file for this gin",
				Required: true,
			},
		},
		Subcommands: commands(),
	}
}
