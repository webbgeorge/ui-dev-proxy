package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/webbgeorge/ui-dev-proxy/commands"
	"github.com/webbgeorge/ui-dev-proxy/file"
)

var Version string

func main() {
	app := cli.NewApp()
	app.Name = "ui-dev-proxy"
	app.Version = Version

	logger := log.New(os.Stdout, "", log.LstdFlags)
	app.Writer = logger.Writer()

	confProvider := file.ConfigProvider()

	app.Commands = []cli.Command{
		commands.StartCommand(logger, confProvider),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
