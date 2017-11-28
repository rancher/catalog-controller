package main

import (
	"os"

	"github.com/rancher/catalog-controller/controller"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var VERSION = "v0.0.0-dev"

func main() {
	app := cli.NewApp()
	app.Name = "catalog-controller"
	app.Version = VERSION
	app.Author = "Rancher Labs, Inc."
	app.Before = func(ctx *cli.Context) error {
		if ctx.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug",
		},
		cli.StringFlag{
			Name:   "config",
			Usage:  "Kube config for accessing kubernetes cluster",
			EnvVar: "KUBECONFIG",
		},
		cli.StringFlag{
			Name:  "cache-root",
			Usage: "Cache root for catalog controller",
		},
		cli.IntFlag{
			Name:  "refresh-interval",
			Usage: "Refresh interval for catalog",
			Value: 60,
		},
	}

	app.Action = controller.Run
	app.Run(os.Args)
}
