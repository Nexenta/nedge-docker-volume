package nedcli

import (
	"github.com/codegangsta/cli"
	"github.com/Nexenta/nedge-docker-volume/nedv/daemon"
)

var (
	BarCmd = cli.Command{
		Name:  "daemon",
		Usage: "daemon related commands",
		Subcommands: []cli.Command{
			BarStartCmd,
		},
	}

	BarStartCmd = cli.Command{
		Name:  "start",
		Usage: "Start the Nedge Docker Daemon: `start [options] NAME`",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "verbose, v",
				Usage: "Enable verbose/debug logging: `[--verbose]`",
			},
			cli.StringFlag{
				Name:  "config, c",
				Usage: "Config file for daemon (default: /opt/nedge/etc/ccow/ned.json): `[--config /opt/nedge/etc/ccow/ned.json]`",
			},
		},
		Action: cmdDaemonStart,
	}
)

func cmdDaemonStart(c *cli.Context) {
	verbose := c.Bool("verbose")
	cfg := c.String("config")
	if cfg == "" {
		cfg = "/opt/nedge/etc/ccow/ned.json"
	}
	daemon.Start(cfg, verbose)
}
