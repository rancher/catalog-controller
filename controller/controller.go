package controller

import (
	"os"
	"os/signal"
	"path"
	"syscall"

	"context"

	"fmt"

	"time"

	"github.com/rancher/catalog-controller/client"
	"github.com/rancher/catalog-controller/manager"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func Run(ctx *cli.Context) error {
	logrus.Infof("Starting catalog controller")
	clientset, err := client.NewClientSetV1(ctx.GlobalString("config"))
	if err != nil {
		return err
	}
	cacheRoot := ctx.GlobalString("cache-root")
	if cacheRoot == "" {
		cacheRoot = path.Join(os.Getenv("HOME"), ".catalog-controller", "cache")
	}
	m := manager.New(clientset, cacheRoot)

	context, cancel := context.WithCancel(context.Background())
	controller := clientset.CatalogClientV1.Catalogs("").Controller()
	controller.AddHandler(m.Sync)
	controller.Start(context, 1)

	interval := ctx.Int("refresh-interval")
	go func() {
		for {
			ticker := time.Tick(time.Duration(interval) * time.Second)
			select {
			case <-ticker:
				catalogs, err := m.GetCatalogs()
				if err != nil {
					logrus.Error(err)
					continue
				}
				for _, catalog := range catalogs {
					controller.Enqueue("", catalog.Name)
				}
			}
		}
	}()

	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	select {
	case <-term:
		logrus.Infof("Received SIGTERM, shutting down")
		os.Exit(0)
	case <-context.Done():
		cancel()
	}
	fmt.Println("exiting")
	return nil
}
