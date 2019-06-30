package main

import (
	"os"

	"github.com/kassybas/tame/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "tame"
	app.Usage = "tame executes targets defined in tame.yaml"
	app.Version = "0.1.4"

	var tameFile string
	var targetName string

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "file, f",
			Value:       "tame.yaml",
			Usage:       "Path of tame source yaml",
			Destination: &tameFile,
		},
		cli.StringSliceFlag{
			Name:  "args, a",
			Usage: "Set arguments for target [arg-name=arg-value]",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.NArg() > 0 {
			targetName = c.Args()[0]
		}
		targetArgs := c.GlobalStringSlice("args")
		cmd.MakeCommand(tameFile, targetName, targetArgs)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
}
