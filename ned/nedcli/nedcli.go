package nedcli

import (
	"fmt"
	"github.com/codegangsta/cli"
)


func NedCmdNotFound(c *cli.Context, command string) {
	fmt.Println(command, " not found ");

}

func NedInitialize(c *cli.Context) error {

	cfgFile := c.GlobalString("config")
	if cfgFile != "" {
		fmt.Println("Found config: ", cfgFile);
	}
	return nil
}

func NewCli(version string) *cli.App {
	app := cli.NewApp()
	app.Name = "ned"
	app.Version = version
	app.Author = "nexentaedge@nexenta.com"
	app.Usage = "CLI for NexentaEdge clusters"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "loglevel",
			Value:  "info",
			Usage:  "Specifies the logging level (debug|warning|error)",
			EnvVar: "LogLevel",
		},
	}
	app.CommandNotFound = NedCmdNotFound
	app.Before = NedInitialize
	app.Commands = []cli.Command{
		DaemonCmd,
		VolumeCmd,
	}
	return app
}
