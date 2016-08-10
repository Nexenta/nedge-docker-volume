package ndvolcli

import (
	"github.com/urfave/cli"
	"github.com/Nexenta/nedge-docker-volume/ndvol/daemon"
)

var (
	DaemonCmd = cli.Command{
		Name:  "daemon",
		Usage: "daemon related commands",
		Subcommands: []cli.Command{
			DaemonStartCmd,
		},
	}

	DaemonStartCmd = cli.Command{
		Name:  "start",
		Usage: "Start the Nedge Docker Daemon: `start [options] NAME`",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "verbose, v",
				Usage: "Enable verbose/debug logging: `[--verbose]`",
			},
			cli.StringFlag{
				Name:  "config, c",
				Usage: "Config file for daemon (default: /etc/nvd/nvd.json): `[--config /etc/nvd/nvd.json]`",
			},
		},
		Action: cmdDaemonStart,
	}
)

func cmdDaemonStart(c *cli.Context) {
	verbose := c.Bool("verbose")
	cfg := c.String("config")
	if cfg == "" {
		cfg = "/opt/nedge/etc/ccow/ndvol.json"
	}
	daemon.Start(cfg, verbose)
}
