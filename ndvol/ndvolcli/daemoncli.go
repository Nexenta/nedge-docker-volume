package ndvolcli

import (
	"github.com/urfave/cli"
	ndvolDaemon "github.com/Nexenta/nedge-docker-volume/ndvol/daemon"
	"github.com/sevlyar/go-daemon"
	log "github.com/Sirupsen/logrus"
	"syscall"
)

var (
	DaemonCmd = cli.Command{
		Name:  "daemon",
		Usage: "daemon related commands",
		Subcommands: []cli.Command{
			DaemonStartCmd,
			DaemonStopCmd,
		},
	}

	DaemonStartCmd = cli.Command{
		Name:  "start",
		Usage: "Start the Nedge Docker Daemon: `start [flags]`",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "verbose, v",
				Usage: "Enable verbose/debug logging: `[--verbose]`",
			},
			cli.StringFlag{
				Name:  "config, c",
				Usage: "Config file for daemon (default: /opt/nedge/etc/ccow/nvd.json): `[--config /opt/nedge/etc/ccow/nvd.json]`",
			},
		},
		Action: cmdDaemonStart,
	}

	DaemonStopCmd = cli.Command{
		Name: "stop",
		Usage: "Stop the Nedge Docker Daemon: `stop",
		Action: cmdDaemonStop,
	}
)

func cmdDaemonStop(c *cli.Context) {
	cntxt := &daemon.Context{
		PidFileName: "/opt/nedge/var/run/ndvol.pid",
		PidFilePerm: 0644,
		LogFileName: "/opt/nedge/var/log/ndvol.log",
		LogFilePerm: 0640,
		Umask:       027,
	}
	d, err := cntxt.Search()
	if err != nil {
		log.Fatalln("Unable to send signal to the daemon:", err)
	}
	d.Signal(syscall.SIGTERM)
}

func cmdDaemonStart(c *cli.Context) {
	cntxt := &daemon.Context{
		PidFileName: "/opt/nedge/var/run/ndvol.pid",
		PidFilePerm: 0644,
		LogFileName: "/opt/nedge/var/log/ndvol.log",
		LogFilePerm: 0640,
		Umask:       027,
	}
	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	defer cntxt.Release()
	if d != nil {
		return
	}

	log.Info("- - - - - - - - - - - - - - -")
	log.Info("Daemon started")
	go DaemonStart(c)

	err = daemon.ServeSignals()
	if err != nil {
		log.Info("Error:", err)
	}
	log.Info("Daemon terminated")
}

func DaemonStart(c *cli.Context) {
	verbose := c.Bool("verbose")
	cfg := c.String("config")
	if cfg == "" {
		cfg = "/opt/nedge/etc/ccow/ndvol.json"
	}
	ndvolDaemon.Start(cfg, verbose)
}
