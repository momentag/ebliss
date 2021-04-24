package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func setGinMode(mode bool) {
	if !mode {
		gin.SetMode(gin.DebugMode)
		log.Info().Str("gin mode", "debug").Send()
	} else {
		gin.SetMode(gin.ReleaseMode)
		log.Info().Str("gin mode", "release").Send()
	}
}

func runServer(port string, addr string, mode bool) error {
	var router = gin.New()
	setGinMode(mode)
	router.POST("/login", Login)
	if err := router.Run(fmt.Sprintf("%v:%v", addr, port)); err != nil {
		log.Error().AnErr("server", err)
		return err
	}
	return nil
}

func authServer() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "starts an http server for authentication",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "port", Aliases: []string{"p"}, Usage: "the port to bind to", DefaultText: "333", Value: "333", HasBeenSet: true},
			&cli.StringFlag{Name: "address", Aliases: []string{"a"}, Usage: "the address to bind for listening", DefaultText: "0.0.0.0", Value: "0.0.0.0", HasBeenSet: true},
			&cli.BoolFlag{Name: "gin", Aliases: []string{"g"}, Usage: "enables production level for gin", DefaultText: "false", Value: false, HasBeenSet: false},
		},
		Action: func(context *cli.Context) error {
			var port = context.String("port")
			var addr = context.String("address")
			var mode = context.Bool("gin")
			log.Info().Str("port", port).Str("addr", addr).Msg("started server")
			return runServer(port, addr, mode)
		},
	}
}
