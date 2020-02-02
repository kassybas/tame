package main

import (
	"os"

	"github.com/kassybas/tame/cmd"
	"github.com/kassybas/tame/internal/helpers"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "tame"
	app.Usage = "tame executes targets defined in tame.yaml"
	app.Version = "0.2.0"

	var tameFile string
	var targetName string

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "file, f",
			Value:       "tame.yaml",
			Usage:       "Path of tame source yaml",
			Destination: &tameFile,
		},
	}

	app.Action = func(c *cli.Context) {
		if c.NArg() > 0 {
			targetName = c.Args()[0]
		}
		targetArgs, err := helpers.ParseCLITargetArgs(c.Args())
		if err != nil {
			logrus.Fatalf("failed to parse target flags:\n\t%s", err)

		}
		cmd.MakeCommand(tameFile, targetName, targetArgs)
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Error(err)
	}
}
