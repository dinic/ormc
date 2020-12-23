package main

import (
	"os"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ormc"
	app.Version = "0.0.1"
	app.Compiled = time.Now()
	app.Copyright = "(C) dinic x.dinic@gmail.com"
	app.Usage = "ormc code|config [-c config path]"

	app.Commands = []cli.Command{
		cmdCode,
		cmdConfig,
	}

	app.Action = func(c *cli.Context) error {
		cli.ShowAppHelp(c)
		return nil
	}

	app.Run(os.Args)
}
